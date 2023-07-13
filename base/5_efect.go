package base

var DotTags = []string{ 
	"DE", // Element
	"DW", // Weight
	"DD", // Dot Delay
}

var RegenTags = []string{ 
	"HPR",  // HP promille portion
}

type Effect struct {
	Time int
	Collision [2]int
	Effects []interface{}
}

func NewEffect() *Effect {
	buffer := Effect{ Time: Epoch() }
	return &buffer
}

func (ef *Effect) Add_Self_MakeDot(dot *Dot) {
	buffer := make(map[string]int)
	buffer["DD"] = CeilRound(1618/dot.Weight()+1)
	buffer["DE"] = ElemIndex[dot.Elem()]
	buffer["DW"] = (*dot)[dot.Elem()]
	(*ef).Effects = append((*ef).Effects, buffer)
}

func (ef *Effect) Add_Self_HPRegen(hp int) {
	buffer := make(map[string]int)
	buffer["HPR"] = hp
	(*ef).Effects = append((*ef).Effects, buffer)
}