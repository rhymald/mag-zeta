package play

import (
	"rhymald/mag-zeta/base"
)

func MakeNPC() *Character {
	var buffer Character 
	buffer.ID = base.Epoch()
	buffer.Body = base.MakeStream()
	buffer.Energy = append(buffer.Energy, base.MakeStream())
	buffer.Life = base.MakeLife()
	return &buffer
}