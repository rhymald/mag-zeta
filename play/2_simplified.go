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
	RXY struct {
		RNow int `json:"RNow"`
		Rotate int `json:"Rotate"`
		XYNow [2]int `json:"XYNow"`
		XYBefore [3][2]int `json:"XYBefore"`
	} `json:"RXY"`
	Look struct {
		Move map[string][2]int `json:"Move"` // how : from
		Cast map[string][2]int `json:"Cast"` // what (tool/fractal/): total ms, left
		Drag map[string]string `json:"Drag"` // what : where - arm[s], shoulder[s], back, belt, neck, leg[s], head
	} `json:"Look"`
}

func (c *Character) Simplify(path [5][2]int, camera [2]int) Simplified {
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
	buffer.TS = make(map[string]int) // (*c).ID
	buffer.TS["Born"] = (*c).TSBorn
	buffer.TS["Atts"] = (*c).TSAtts
	c.Unlock()
	buffer.RXY.XYNow = [2]int{ camera[0]-path[1][0], camera[1]-path[1][1] }
	buffer.RXY.RNow = path[0][0]
	buffer.RXY.Rotate = path[0][1]
	for i, each := range path[2:5] { buffer.RXY.XYBefore[i] = [2]int{ camera[0]-each[0], camera[1]-each[1] } }
	// immitation:
	barrier, penalty := base.CeilRound(100*base.Rand()), base.FloorRound(100*base.Rand())
	buffer.Wound = penalty
	buffer.Barrier = barrier
	buffer.Attune = "TBD"
	return buffer
}