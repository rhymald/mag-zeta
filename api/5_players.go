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
		(*world).ByID[player.GetID()] = player
		c.IndentedJSON(200, "Successfully logged in")
	} else {
		c.AbortWithError(500, errors.New("Invalid player character"))
	}
	world.Unlock()
	spanResponse.End()
	
	go func(){ playerLiveAlive(player, (*c).Request.Context()) }()
	// select {}
}

func playerRegen(c *play.Character, ctx context.Context) {
	// if c.IsNPC() == false {
		// _, span := tracer.Start(ctx, "generate-dot")
		// defer span.End()
		picker := base.EpochNS() % len((*c).Energy)
		// span.SetAttributes(attribute.Int("ByStream", picker))
		_, dot := c.GetDotFrom(picker)
		// span.SetAttributes(attribute.Int("DotIdx", idx))
		// span.AddEvent(dot.ToStr())
		base.Wait(256*dot.Weight())
		hp := 3
		c.Life.HealDamage(hp)
		(*c).ID["Life"] = base.Epoch()
		// span.SetAttributes(attribute.Int("HPGain", hp))
	// } else {
	// 	// _, span := tracer.Start(ctx, "npc-regeneration")
	// 	// defer span.End()
	// 	base.Wait(256)
	// 	hp := 1
	// 	c.Life.HealDamage(hp)
	// 	(*c).ID["Life"] = base.Epoch()
	// 	// span.SetAttributes(attribute.Int("HPGain", hp))
	// }
}

func playerLiveAlive(c *play.Character, ctx context.Context) {
	// ctx2, span := tracer.Start(ctx, "lifecycle-regeneration")
	// defer span.End()
	for {
		c.Lock()
		if c.Life.Dead() { return } //span.AddEvent("Character died")
		energyFull := len((*c).Pool) < base.ChancedRound((*(*c).Atts).Capacity)
		c.Unlock()
		if energyFull { base.Wait(4096) ; continue } //span.AddEvent("Energy full, wait") ;
		c.Lock()
		playerRegen(&*c, ctx) //2)
		c.Unlock()
	}
}
// + regen
// + potion(s)
// + move
// + jinx[e], punch[p] 