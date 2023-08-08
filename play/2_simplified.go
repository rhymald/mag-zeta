package play 

import (
	"rhymald/mag-zeta/base"
	"fmt"
) 

type Simplified struct {
	// E string `json:"E"`
	Name string `json:"Name"`
	ID string `json:"ID"`
	TS map[string]int `json:"TS"` 
	// health
	HP int `json:"HP"`
	// Barrier int `json:"Barrier"`
	// Wound int `json:"Wound"`
	// elem
	PWR int `json:"PWR"`
	// xyz
	RXY struct {
		RNow int `json:"RNow"`
		RAdd int `json:"RAdd"`
		XYNow [2]int `json:"XYNow"`
		XYOld [3][2]int `json:"XYOld"`
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
		buffer.PWR = -base.ChancedRound((*(*c).Atts).Capacity) 
	} else { 
		// buffer.ID = "Player"
		buffer.PWR = len((*c).Pool)
	}
	buffer.ID = c.GetID()
	buffer.TS = make(map[string]int) // (*c).ID
	buffer.TS["Born"] = (*c).TSBorn
	buffer.TS["Atts"] = (*c).TSAtts
	body, elem := (*c).Body, (*c).Energy[0]
	eb, ee := body.Elem(), elem.Elem()
	if eb == base.ElemList[0] { eb = "" }
	if ee == base.PhysList[0] { if eb == "" { ee = "🧿" } else { ee = ""}}
	c.Unlock()
	buffer.RXY.XYNow = [2]int{ camera[0]-path[1][0], camera[1]-path[1][1] }
	buffer.RXY.RNow = path[0][0]
	buffer.RXY.RAdd = path[0][1]
	for i, each := range path[2:5] { buffer.RXY.XYOld[i] = [2]int{ camera[0]-each[0], camera[1]-each[1] } }
	if npc { 
		// buffer.E = fmt.Sprintf("%s%s", eb, ee)
		buffer.Name = fmt.Sprintf("[%s%s] Training dummy", eb, ee)
	} else { 
		// buffer.E = fmt.Sprintf("%s", eb) 
		buffer.Name = "Some player"
	}
	// immitation:
	// barrier, penalty := base.CeilRound(100*base.Rand()), base.FloorRound(100*base.Rand())
	// buffer.Wound = penalty
	// buffer.Barrier = barrier
	return buffer
}