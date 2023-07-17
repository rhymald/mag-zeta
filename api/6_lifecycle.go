package api

import (
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/base"
	"context"
	// "go.opentelemetry.io/otel/trace"
	"math"
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

func Lifecycle_Regenerate(st *play.State, ctx context.Context) {
	isNPC := (*st).Current.IsNPC()
	if !isNPC { for {
		(*st).Current.Lock()
		if len((*(*st).Current).Pool) >= base.ChancedRound((*(*(*st).Current).Atts).Capacity) { (*st).Current. Unlock() ; base.Wait(4236) ; break }
		(*st).Current.Unlock()
		_, span := tracer.Start(ctx, fmt.Sprintf("player-%s-regen", (*st).Current.GetID()))
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
		span.AddEvent(fmt.Sprintf("%d|+%d[%s]|+HP[%d]", picker, (*effect).Time, dot.ToStr(), 8))
		// fmt.Println((*st).Effects)
		st.Unlock()
		base.Wait(1618/dot.Weight()+1)
	}} else { for {
		_, span := tracer.Start(ctx, fmt.Sprintf("npc-%s-regen", (*st).Current.GetID()))
		defer span.End()
		if (*st).Current.Life.Wounded() { span.AddEvent("NPC died") ; return }
		effect := base.NewEffect()
		effect.Add_Self_HPRegen(32)
		st.Lock()
		(*st).Effects[(*effect).Time] = effect
		// fmt.Println((*st).Effects)
		st.Unlock()
		span.AddEvent(fmt.Sprintf("%d|+%d[%s]|+HP[%d]", -1, (*effect).Time, "none", 32))
		base.Wait(4236)
	}}
}

func Lifecycle_EffectConsumer(st *play.State, ctx context.Context) {
	pause, now := 1618, base.Epoch()
	first, sum, counter := 0, 0, 0
	prefix := "player" ; if (*st).Current.IsNPC() { prefix = "npc" }
	ctxLifeCycle, span := tracer.Start(ctx, fmt.Sprintf("%s-%s-regen", prefix, (*st).Current.GetID()))
	defer span.End()
	for {
		if len((*st).Effects) == 0 { base.Wait(float64(pause)) ; continue }
		fmt.Println("Before:", (*st).Effects, sum)

		// step 1 read to limit
		_, spanReader := tracer.Start(ctxLifeCycle, "take-effects")
		buffer := make(map[int]*base.Effect)
		st.Lock()
		for ts, effect := range (*st).Effects {
			if len(buffer) == 0 { first = ts }
			if ts-first > pause { continue }
			buffer[ts] = effect
			counter++ ; sum += ts - first
			if sum > counter * pause { break }
		}
		spanReader.AddEvent(fmt.Sprintf("Effects: { read: %d, total: %d }", counter, len((*st).Effects)))
		st.Unlock()
		spanReader.End()

		// step 2 sum and distribute
		// instant, conditions, delayed := []base.Effect{}, []base.Effect{}, []base.Effect{}
		// timerInst, timerCond, timerDel := 0, 0, 0
		// for 
	
		// step 3 consume
	
		// step 4 clean read from queue
		_, spanDeleter := tracer.Start(ctxLifeCycle, "take-effects")
		st.Lock() 
		for ts, _ := range buffer { delete((*st).Effects, ts) }
		spanDeleter.AddEvent(fmt.Sprintf("Effects: { read: %d, total: %d }", counter, len((*st).Effects)))
		st.Unlock()
		spanDeleter.End()
		
		// step 5 redirect back leftovers
	
		// end
		delay := float64(pause) / ( 1/math.Phi + math.Log2(1+math.Abs(float64(now-first))) )
		fmt.Println("After:", (*st).Effects, sum)
		span.AddEvent(fmt.Sprintf("WaitFor: %0.3fms", delay ))
		base.Wait( delay )
	}
}