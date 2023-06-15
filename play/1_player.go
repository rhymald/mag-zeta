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

func (c *Character) GetDotFrom(strIndex int) (int, *base.Dot) {
	c.Lock()
	buffer := (*c).Pool
	stream := (*c).Energy[strIndex]
	index := base.EpochNS()
	for { _, ok := buffer[index] ; if ok { index++} else { break } }
	buffer[index] = stream.MakeDot()
	(*c).Pool = buffer
	(*c).ID["Pool"] = base.EpochNS()
	c.Unlock()
	return index, buffer[index]
}

func (c *Character) BurnDot() (int, *base.Dot) {
	if c.IsNPC() { str := (*c).Energy[0] ; base.Wait(5) ; return base.EpochNS(), str.MakeDot() }
	dot, tstamp := &base.Dot{}, -1
	c.Lock()
	buffer := (*c).Pool
	min := base.Epoch()
	for ts := range buffer {
		// tstamp, dot = ts, (*c).Pool[ts]
		// delete(buffer, ts)
		// break
		if ts < min { min = ts }
	}
	if len(buffer) != 0 {
		tstamp, dot = min, (*c).Pool[min]
		delete(buffer, min)
	}
	(*c).Pool = buffer
	c.Unlock()
	return tstamp, dot
}