package play

import (
	"rhymald/mag-zeta/base"
)

func MakePlayer() *Character {
	var buffer Character 
	buffer.ID = make(map[string]int)
	buffer.ID["Born"] = base.Epoch()	
	buffer.Body = base.MakeStream()
	for x:=0; x<LuckyBorn(buffer.ID["Born"]); x++ { buffer.Energy = append(buffer.Energy, base.MakeStream()) }
	buffer.Life = base.MakeLife()
	buffer.Pool = make(map[int]*base.Dot)
	return &buffer
}

func (c *Character) GetDotFrom(strIndex int) {
	base.Wait(64)
	c.Lock()
	buffer := (*c).Pool
	stream := (*c).Energy[strIndex]
	buffer[base.Epoch()] = stream.MakeDot()
	(*c).Pool = buffer
	c.Unlock()
	base.Wait(64)
}

func (c *Character) BurnDot() (int, *base.Dot) {
	dot, tstamp := &base.Dot{}, -1
	c.Lock()
	buffer := (*c).Pool
	for ts, _ := range buffer {
		tstamp, dot = ts, (*c).Pool[ts]
		delete(buffer, ts)
		break
	}
	(*c).Pool = buffer
	c.Unlock()
	return tstamp, dot
}