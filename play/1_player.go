package play

import (
	"rhymald/mag-zeta/base"
)

func MakePlayer() *Character {
	var buffer Character 
	buffer.ID = base.Epoch()
	buffer.Body = base.MakeStream()
	for x:=0; x<LuckyBorn(buffer.ID); x++ { buffer.Energy = append(buffer.Energy, base.MakeStream()) }
	buffer.Life = base.MakeLife()
	return &buffer
}