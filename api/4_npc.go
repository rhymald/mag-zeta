package api

import (
	"errors"
	"rhymald/mag-zeta/play"
	// "rhymald/mag-zeta/base"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	// "rhymald/mag-zeta/base"
	// "go.opentelemetry.io/otel/trace"
	// "fmt"
	"math"
	// "rhymald/mag-zeta/connect"
	"fmt"
)

func newFoe(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "spawn-foe")
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
	world.Lock()
	state := foe.NewState()
	if err == nil {
		(*world).ByID[foe.GetID()] = state
		c.IndentedJSON(200, struct{ ID string }{ ID: foe.GetID() })
	} else {
		c.AbortWithError(500, errors.New("Invalid foe character"))
	}
	world.Unlock()
	spanResponse.End()

	go func(){ Lifecycle_Regenerate(state, (*c).Request.Context()) }()
	go func(){ Lifecycle_EffectConsumer(state, (*c).Request.Context()) }()
	go func(){ for {state.Move(GridCache)} ; state.Turn(1/math.Phi/math.Phi, GridCache) }()
	// go func(){ for {state.Turn(1/math.Phi/math.Phi * float64(base.Epoch()%3-1), GridCache)} }()
	// go func(){ base.Wait(30000) ; state.Lock() ; connect.WriteTrace((*world).Writer, foe.GetID(), &(*state).Trace) ; state.Unlock() }()
}

// func npcRegen(hps *base.Life, ids *map[string]int, span *trace.Span) {
// 	hp := 32
// 	hps.HealDamage(hp)
// 	(*ids)["Life"] = base.Epoch()
// 	(*span).AddEvent(fmt.Sprintf("0|+0[none|-1]+HP|%+d", hp))
// }
