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