package play

import (
	"rhymald/mag-zeta/base"
)

func MakeNPC() *Character {
	var buffer Character 
	// buffer.ID = make(map[string]int)
	buffer.TSBorn = base.Epoch()
	buffer.Body = base.MakeStream(base.PhysList[2])
	buffer.Energy = append(buffer.Energy, base.MakeStream(base.ElemList[0]))
	buffer.Life = base.MakeLife()
	return &buffer
}