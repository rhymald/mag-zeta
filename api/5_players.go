package api

import (
	"errors"
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/base"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"fmt"
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
	if err == nil {
		(*world).ByID[player.GetID()] = player.NewState()
		c.IndentedJSON(200, "Successfully logged in")
	} else {
		c.AbortWithError(500, errors.New("Invalid player character"))
	}
	world.Unlock()
	spanResponse.End()
	
	go func(){ charLiveAlive(player, (*c).Request.Context()) }()
}

func playerRegen(hps *base.Life, pool *map[int]*base.Dot, ids *map[string]int, energy *[]*base.Stream, span *trace.Span) float64 {
	picker := base.EpochNS() % len(*energy)
	stream := (*energy)[picker]
	idx, dot := play.GetDotFrom(pool, stream, ids)
	hp := 8
	hps.HealDamage(hp)
	(*ids)["Life"] = base.Epoch()
	(*span).AddEvent(fmt.Sprintf("%d|+%d[%s]|+HP[%d]", picker, idx, dot.ToStr(), hp))
	return 1000*dot.Weight()
}
