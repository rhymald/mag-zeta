package play

import (
	"rhymald/mag-zeta/base"
)

func MakePlayer() *Character {
	var buffer Character 
	// buffer.ID = make(map[string]int)
	buffer.TSBorn = base.Epoch()	
	buffer.Body = base.MakeStream(base.PhysList[1])
	for x:=0; x<LuckyBorn(buffer.TSBorn); x++ { buffer.Energy = append(buffer.Energy, base.MakeStream(base.ElemList[0])) }
	buffer.Life = base.MakeLife()
	buffer.Pool = make(map[int]*base.Dot)
	return &buffer
}

// func GetDotFrom(pool *map[int]*base.Dot, stream *base.Stream, ids *map[string]int) (int, *base.Dot) {
// 	index := base.Epoch()
// 	for { _, ok := (*pool)[index] ; if ok { index++ } else { break } }
// 	(*pool)[index] = stream.MakeDot()
// 	(*ids)["Pool"] = base.Epoch()
// 	return index, (*pool)[index]
// }

func (c *Character) BurnDot() (int, *base.Dot) {
	if c.IsNPC() { str := (*c).Energy[0] ; base.Wait(5) ; return base.Epoch(), str.MakeDot() }
	dot, tstamp := &base.Dot{}, -1
	c.Lock()
	buffer := (*c).Pool
	min := base.Epoch()
	for ts := range buffer { if ts < min { min = ts } }
	if len(buffer) != 0 {
		tstamp, dot = min, (*c).Pool[min]
		delete(buffer, min)
	}
	(*c).Pool = buffer
	c.Unlock()
	return tstamp, dot
}