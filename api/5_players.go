package api

import (
	"errors"
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/base"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"context"
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
		(*world).Players = append((*world).Players, player)
		c.IndentedJSON(200, "Successfully logged in")
	} else {
		c.AbortWithError(500, errors.New("Invalid player character"))
	}
	spanResponse.End()
	
	go func(){ playerLiveAlive((*world).Players[len((*world).Players)-1], (*c).Request.Context()) }()
	world.Unlock()
}

func playerRegen(c *play.Character, ctx context.Context) {
	if c.IsNPC() {
		_, span := tracer.Start(ctx, "generate-dot")
		defer span.End()
		picker := base.EpochNS() % len((*c).Energy)
		span.SetAttributes(attribute.Int("ByStream", picker))
		idx, dot := c.GetDotFrom(picker)
		span.SetAttributes(attribute.Int("DotIdx", idx))
		span.AddEvent(dot.ToStr())
		base.Wait(256*dot.Weight())
		hp := 3
		c.Life.HealDamage(hp)
		(*c).ID["Life"] = base.Epoch()
		span.SetAttributes(attribute.Int("HPGain", hp))
	} else {
		_, span := tracer.Start(ctx, "npc-regeneration")
		defer span.End()
		base.Wait(256)
		hp := 1
		c.Life.HealDamage(hp)
		(*c).ID["Life"] = base.Epoch()
		span.SetAttributes(attribute.Int("HPGain", hp))
	}
}

func playerLiveAlive(c *play.Character, ctx context.Context) {
	ctx2, span := tracer.Start(ctx, "lifecycle-regeneration")
	defer span.End()
	for {
		if c.Life.Dead() { span.AddEvent("Character died") ; return }
		energyFull := len((*c).Pool) < base.ChancedRound((*(*c).Atts).Capacity)
		if energyFull { span.AddEvent("Energy full, wait") ; base.Wait(4096) ; continue }
		playerRegen(&*c, ctx2)
	}
}
// + regen
// + potion(s)
// + move
// + jinx[e], punch[p] 