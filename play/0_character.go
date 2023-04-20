package play

import (
	"rhymald/mag-zeta/base"
	"errors"
)

type Attributes struct { 
	Vitality float64 
	Agility float64
	Capacity float64
	Resistance map[string]float64
}

type Character struct {
	ID int
	// basics
	Body *base.Stream
	Energy []*base.Stream
	// consumables
	Life *base.Life
	Pool []*base.Dot
	// recalculateable stats
	Atts *Attributes
}

func LuckyBorn(time int) int { if time%10 == 0 {return 2} else if time%10 == 9 {return 5} else if time%10 < 5 {return 3} else {return 4} ; return 0}
func (c *Character) IsNPC() bool { return len((*cr).Energy) <= 1 }

func (c *Character) CalculateAttributes() error {
	if (*c).ID < 1000000 { return errors.New("Character Attributes: Empty character ID.") }
	if len((*c).Energy) == 0 { return errors.New("Character Attributes: Empty character streams.") }
	var buffer Attributes
	buffer.Vitality = (*c).Body.Dot() * 10
	buffer.Agility = (*c).Body.Mean() * 0.7
	buffer.Resistance = make(map[string]float64)
	mod := float64(6 - LuckyBorn((*c).ID))
	for _, each := range (*c).Energy { 
		buffer.Resistance[each.Elem()] += each.Mean()
		buffer.Capacity += each.Len() * mod
	}
	(*c).Atts = &buffer
	return nil
}