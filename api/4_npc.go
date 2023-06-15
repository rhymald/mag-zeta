package api

import (
	"errors"
	"rhymald/mag-zeta/play"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

func newFoe(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "spawn-foe")
	defer span.End()

	_, spanGenerate := tracer.Start(ctx, "generating-basic-stats")
	foe := play.MakeNPC()
	spanGenerate.AddEvent("Character generated with ID: 123456789-1234-1-1234567")
	span.SetAttributes(attribute.String("CharacterID","123456789-1234-1-1234567"))
	spanGenerate.End()

	_, spanCalculate := tracer.Start(ctx, "calculating-attributes-from-basic")
	err := foe.CalculateAttributes()
	if err != nil { spanCalculate.RecordError(err) }
	spanCalculate.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	world.Lock()
	if err == nil {
		(*world).ByID[foe.GetID()] = foe
		c.IndentedJSON(200, "Successfully spawned")
	} else {
		c.AbortWithError(500, errors.New("Invalid foe character"))
	}
	world.Unlock()
	spanResponse.End()

	go func(){ playerLiveAlive(foe, (*c).Request.Context()) }()
	// select {}
}
