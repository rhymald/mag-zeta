package api

import (
	"errors"
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/base"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/trace"
	// "fmt"
	"math"
	"rhymald/mag-zeta/connect"
)

func newPlayer(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "login-player")
	defer span.End()

	_, spanGenerate := tracer.Start(ctx, "generating-basic-stats")
	player := play.MakePlayer()
	spanGenerate.AddEvent("Character generated with ID: 123456789-1234-1-1234567")
	span.SetAttributes(attribute.String("CharacterID","123456789-1234-1-1234567"))
	spanGenerate.End()

	_, spanCalculate := tracer.Start(ctx, "calculating-attributes-from-basic")
	err := player.CalculateAttributes()
	if err != nil { spanCalculate.RecordError(err) }
	spanCalculate.End()
	
	_, spanResponse := tracer.Start(ctx, "responding")
	world.Lock()
	state := player.NewState()
	if err == nil {
		id := player.GetID()
		(*world).ByID[id] = state
		c.IndentedJSON(200, struct{ ID string }{ ID: player.GetID() })
	} else {
		c.AbortWithError(500, errors.New("Invalid player character"))
	}
	world.Unlock()
	spanResponse.End()
	
	go func(){ Lifecycle_Regenerate(state, (*c).Request.Context()) }()
	go func(){ Lifecycle_EffectConsumer(state, (*c).Request.Context()) }()
	go func(){ for x:=0 ; x<10 ; x++ {state.Move(GridCache)} }()
	go func(){ for x:=0 ; x<25 ; x++ {state.Turn(1/math.Phi/math.Phi * float64(base.Epoch()%3-1), GridCache)} }()
	go func(){ base.Wait(30000) ; connect.WritePosition((*world).Writer, player.GetID(), &(*state).Trace) }()
}
