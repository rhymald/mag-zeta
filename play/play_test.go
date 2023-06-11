package play

import "testing"
import "fmt"

func Test_Characters(t *testing.T) {
	npc := MakeNPC()
	player := MakePlayer()
	_,_ = npc.CalculateAttributes(), player.CalculateAttributes()
	t.Logf("NPC Life: %+v", (*npc).Life)
	t.Logf("NPC Body: %+v", (*npc).Body)
	t.Logf("NPC Nrgy: %+v", (*npc).Energy[0])
	t.Logf("NPC Atts: %+v", (*npc).Atts)
	t.Logf("NPC Pool: %+v", (*npc).Pool)
	t.Logf("--------------------------------------------------")
	t.Logf("Player Life: %+v", (*player).Life)
	t.Logf("Player Body: %+v", (*player).Body)
	energy := "Player Nrgy:"
	for _, each := range (*player).Energy { energy = fmt.Sprintf("%s %+v", energy, each) }
	t.Logf(energy)
	t.Logf("Player Atts: %+v", (*player).Atts)
	t.Logf("Player Pool: %+v", (*player).Pool)
}

func Test_CharCache_Life(t *testing.T) {
	npc := MakeNPC()
	player := MakePlayer()
	_,_ = npc.CalculateAttributes(), player.CalculateAttributes()
	npcCache, playerCache := npc.NewState(), player.NewState()
	t.Logf("--------------------------------------------------")
	t.Logf("Player Cache: %+v", (*playerCache))
	t.Logf("Player Current: %+v", *playerCache.Current.Life)
	t.Logf("NPC Cache: %+v", (*npcCache))
	t.Logf("NPC Current: %+v", *npcCache.Current.Life)
	(*player).Life.Rate = 100
	(*npc).Life.Rate = 900
	npcCache.UpdLife() ; playerCache.UpdLife()
	t.Logf("--------------------------------------------------")
	t.Logf("Player Later: %+v", (*playerCache))
	t.Logf("Player Current: %+v", *playerCache.Current.Life)
	t.Logf("NPC Cache: %+v", (*npcCache))
	t.Logf("NPC Current: %+v", *npcCache.Current.Life)
}

func Test_Dots(t *testing.T) {
	npc := MakeNPC()
	player := MakePlayer()
	_,_ = npc.CalculateAttributes(), player.CalculateAttributes()
	t.Logf("--------------------------------------------------")
	t.Logf("Player: %+v", (*player).Pool)
	for x:=0 ; x<12 ; x++ { 
		player.GetDotFrom( x % len( (*player).Energy ) ) 
		t.Logf("    +Gain: %+v", (*player).Pool)
		if x%3 != 0 { ts, dot := player.BurnDot()  ; t.Logf("    -Burn: %d:%+v %+v", ts, dot, (*player).Pool) }
	}
	t.Logf("--------------------------------------------------")
	t.Logf("NPC: %+v", (*npc).Pool)
	for x:=0 ; x<5 ; x++ {
		ts, dot := npc.BurnDot()
		t.Logf("    -Burn: %d:%+v", ts, dot)
	}
}