package api

import (
	"github.com/gin-gonic/gin"
	"errors"
	"rhymald/mag-zeta/play"
	// "rhymald/mag-zeta/base"
	"go.opentelemetry.io/otel/attribute"
	// "rhymald/mag-zeta/base"
	// "go.opentelemetry.io/otel/trace"
	// "fmt"
	// "math"
	// "rhymald/mag-zeta/connect"
	"fmt"
)

func showGrid(c *gin.Context) { 
	c.IndentedJSON(200, *world) 
}

func showState(c *gin.Context) { 
	id := c.Param("id")
	c.IndentedJSON(200, (*world).ByID[ id ]) 
}

func newFoe(c *gin.Context) { 
	ginContext := (*c).Request.Context()
	ctx, span := tracer.Start(ginContext, "spawn-foe")
	defer span.End()

	_, spanGenerate := tracer.Start(ctx, "generating-basic-stats")
	foe := play.MakeNPC()
	spanGenerate.AddEvent(fmt.Sprintf("Character generated with ID: %s", foe.GetID()))
	span.SetAttributes(attribute.String("CharacterID",foe.GetID()))
	spanGenerate.End()

	_, spanCalculate := tracer.Start(ctx, "calculating-attributes-from-basic")
	err := foe.CalculateAttributes()
	if err != nil { spanCalculate.RecordError(err) }
	spanCalculate.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	state := foe.NewState()
	if err == nil {
		world.Lock()
		renew := (*world).ByID
		renew[foe.GetID()] = state
		(*world).ByID = renew
		world.Unlock()
		c.IndentedJSON(200, struct{ ID string }{ ID: foe.GetID() })
	} else {
		c.AbortWithError(500, errors.New("Invalid foe character"))
	}
	spanResponse.End()

	go func(){ Lifecycle_Regenerate(state, ginContext) }()
	go func(){ Lifecycle_EffectConsumer(state, ginContext) }()
	go func(){ for {
		state.Move(0.11, false, GridCache)
		// state.Turn(0.11, GridCache)
	}}()
	// go func(){ for {state.Turn(1/math.Phi/math.Phi * float64(base.Epoch()%3-1), GridCache)} }()
	// go func(){ base.Wait(30000) ; state.Lock() ; connect.WriteTrace((*world).Writer, foe.GetID(), &(*state).Trace) ; state.Unlock() }()
}
