package api

import (
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/base"
	"context"
	// "go.opentelemetry.io/otel/trace"
)

func charLiveAlive(c *play.Character, ctx context.Context) {
	_, span := tracer.Start(ctx, "lifecycle")
	defer span.End()
	if c.IsNPC() {
		for {
			c.Lock()
			if c.Life.Wounded() { c.Unlock() ; span.AddEvent("Character died") ; return }
			npcRegen((*c).Life, &(*c).ID, &span)
			c.Unlock()
			base.Wait(4096)
		}
	} else {
		for {
			wait := 4096.0
			c.Lock()
			if c.Life.Dead() { span.AddEvent("Character died") ; return }
			energyFull := len((*c).Pool) >= base.ChancedRound((*(*c).Atts).Capacity)
			if energyFull { 
				c.Unlock()
				span.AddEvent("Energy full, wait")
				span.End()
			} else {
				wait = playerRegen((*c).Life, &(*c).Pool, &(*c).ID, &(*c).Energy, &span)
				c.Unlock()
			}
			base.Wait(wait)
		}
	}
}
// [v] regen
// + potion(s)
// + move
// + jinx[e], punch[p] 