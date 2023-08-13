package api

import (
	"errors"
	"rhymald/mag-zeta/play"
	// "rhymald/mag-zeta/base"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/trace"
	// "fmt"
	// "math"
	// "rhymald/mag-zeta/connect"
	"fmt"
)

func newPlayer(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "login-player")
	defer span.End()

	_, spanGenerate := tracer.Start(ctx, "generating-basic-stats")
	player := play.MakePlayer()
	spanGenerate.AddEvent(fmt.Sprintf("Character generated with ID: %s", player.GetID()))
	span.SetAttributes(attribute.String("CharacterID",player.GetID()))
	spanGenerate.End()

	_, spanCalculate := tracer.Start(ctx, "calculating-attributes-from-basic")
	err := player.CalculateAttributes()
	if err != nil { spanCalculate.RecordError(err) }
	spanCalculate.End()
	
	_, spanResponse := tracer.Start(ctx, "responding")
	state := player.NewState()
	if err == nil {
		id := player.GetID()
		world.Lock()
		renew := (*world).ByID
		renew[id] = state
		(*world).ByID = renew
		world.Unlock()
		c.IndentedJSON(200, struct{ ID string }{ ID: player.GetID() })
	} else {
		c.AbortWithError(500, errors.New("Invalid player character"))
	}
	spanResponse.End()
	
	go func(){ Lifecycle_Regenerate(state, (*c).Request.Context()) }()
	go func(){ Lifecycle_EffectConsumer(state, (*c).Request.Context()) }()
	go func(){ for {
		state.Move(1.0/8, true, GridCache)
		// state.Turn(0.09, GridCache)
	}}()
	// go func(){ for {state.Turn(1/math.Phi/math.Phi * float64(base.Epoch()%3-1), GridCache)} }()
	// go func(){ base.Wait(30000) ; state.Lock() ; connect.WriteTrace((*world).Writer, player.GetID(), &(*state).Trace) ; state.Unlock() }()
}
