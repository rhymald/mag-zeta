package api

import (
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/base"
	"context"
	// "go.opentelemetry.io/otel/trace"
	"fmt"
)

// func charLiveAlive(c *play.Character, ctx context.Context) {
// 	_, span := tracer.Start(ctx, "lifecycle")
// 	defer span.End()
// 	if c.IsNPC() {
// 		for {
// 			c.Lock()
// 			if c.Life.Wounded() { c.Unlock() ; span.AddEvent("Character died") ; return }
// 			npcRegen((*c).Life, &(*c).ID, &span)
// 			c.Unlock()
// 			base.Wait(4096)
// 		}
// 	} else {
// 		for {
// 			wait := 4096.0
// 			c.Lock()
// 			if c.Life.Dead() { span.AddEvent("Character died") ; return }
// 			energyFull := len((*c).Pool) >= base.ChancedRound((*(*c).Atts).Capacity)
// 			if energyFull { 
// 				c.Unlock()
// 				span.AddEvent("Energy full, wait")
// 				span.End()
// 			} else {
// 				wait = playerRegen((*c).Life, &(*c).Pool, &(*c).ID, &(*c).Energy, &span)
// 				c.Unlock()
// 			}
// 			base.Wait(wait)
// 		}
// 	}
// }
// [v] regen
// + potion(s)
// + move
// + jinx[e], punch[p] 
// func playerRegen(hps *base.Life, pool *map[int]*base.Dot, ids *map[string]int, energy *[]*base.Stream, span *trace.Span) float64 {
// 	picker := base.EpochNS() % len(*energy)
// 	stream := (*energy)[picker]
// 	idx, dot := play.GetDotFrom(pool, stream, ids) // consumes
// 	hp := 8
// 	hps.HealDamage(hp)
// 	(*ids)["Life"] = base.Epoch()
// 	(*span).AddEvent(fmt.Sprintf("%d|+%d[%s]|+HP[%d]", picker, idx, dot.ToStr(), hp))
// 	return 1000*dot.Weight()
// }

func Regenerate(st *play.State, ctx context.Context) {
	isNPC := (*st).Current.IsNPC()
	if !isNPC { for {
		// if len((*(*st).Current).Pool) < base.ChancedRound((*(*st).Current))
		_, span := tracer.Start(ctx, "player-regen")
		defer span.End()
		effect := base.NewEffect()
		(*st).Current.Lock()
		if (*st).Current.Life.Dead() { (*st).Current.Unlock() ; span.AddEvent("Player died") ; return }
		picker := base.EpochNS() % len((*(*st).Current).Energy)
		stream := (*(*st).Current).Energy[picker]
		(*st).Current.Unlock()
		dot := stream.MakeDot()
		effect.Add_Self_MakeDot(dot)
		effect.Add_Self_HPRegen(8)
		st.Lock()
		(*st).Effects[(*effect).Time] = effect
		fmt.Println((*st).Effects)
		st.Unlock()
		base.Wait(1618/dot.Weight()+1)
	}} else { for {
		_, span := tracer.Start(ctx, "npc-regen")
		defer span.End()
		if (*st).Current.Life.Wounded() { span.AddEvent("NPC died") ; return }
		effect := base.NewEffect()
		effect.Add_Self_HPRegen(32)
		st.Lock()
		(*st).Effects[(*effect).Time] = effect
		fmt.Println((*st).Effects)
		st.Unlock()
		base.Wait(4236)
	}}
}