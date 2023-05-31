package play 

import "rhymald/mag-zeta/base"

type Simplified struct {
	ID string `json:"ID"`
	// health
	HP int `json:"HP"`
	Barrier int `json:"Barrier"`
	Wound int `json:"Wound"`
	// elem
	Attune string `json:"Attune"`
	Power int `json:"Power"`
	// xyz
	XY [2]int `json:"XY"`
	Direction int `json:"Direction"`
}

func (c *Character) Simplify() Simplified {
	var buffer Simplified
	c.Lock()
	npc := c.IsNPC()
	buffer.HP = (*c).Life.Rate
	if npc { 
		buffer.ID = "Dummy"
		buffer.Power = base.ChancedRound((*(*c).Atts).Capacity) 
	} else { 
		buffer.ID = "Player"
		buffer.Power = len((*c).Pool)
	}
	c.Unlock()
	// immitation:
	barrier, penalty := base.CeilRound(100*base.Rand()), base.FloorRound(100*base.Rand())
	buffer.Wound = penalty
	buffer.Barrier = barrier
	buffer.Attune = "TBD"
	buffer.XY = [2]int{ base.CeilRound(200*base.Rand()-100), base.CeilRound(200*base.Rand()-100) }
	buffer.Direction = base.FloorRound(2000*base.Rand()-1000)
	return buffer
}