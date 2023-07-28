package play 

import "rhymald/mag-zeta/base"

type Simplified struct {
	ID string `json:"ID"`
	TS map[string]int `json:"TS"` 
	// health
	HP int `json:"HP"`
	Barrier int `json:"Barrier"`
	Wound int `json:"Wound"`
	// elem
	Attune string `json:"Attune"`
	Power int `json:"Power"`
	// xyz
	// XY [2]int `json:"XY"`
	// Dir int `json:"Direction"`
}

func (c *Character) Simplify() Simplified {
	var buffer Simplified
	c.Lock()
	npc := c.IsNPC()
	buffer.HP = (*c).Life.Rate
	if npc { 
		// buffer.ID = "Dummy"
		buffer.Power = base.ChancedRound((*(*c).Atts).Capacity) 
	} else { 
		// buffer.ID = "Player"
		buffer.Power = len((*c).Pool)
	}
	buffer.ID = c.GetID()
	buffer.TS = (*c).ID
	c.Unlock()
	// buffer.XY = xy
	// buffer.Dir = d
	// immitation:
	barrier, penalty := base.CeilRound(100*base.Rand()), base.FloorRound(100*base.Rand())
	buffer.Wound = penalty
	buffer.Barrier = barrier
	buffer.Attune = "TBD"
	return buffer
}