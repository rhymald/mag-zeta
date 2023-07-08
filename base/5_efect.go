package base

type Effect_Dot map[string]int
var DotTags = []string{ "HP", "DE", "DW", "DT" }

type Effect_DMG map[string]int
var DMGTags = []string{ "HP", "DMG" }

type Effect struct {
	Time int
	Collision [2]int
	Effects []interface{}
}

func NewEffect() *Effect {
	buffer := Effect{ Time: Epoch() }
	return &buffer
}

func (ef *Effect) MakeDot(dot *Dot) {
	buffer := make(map[string]int)
	buffer["DT"] = CeilRound(1000/dot.Weight()+1)
	buffer["DE"] = ElemIndex[dot.Elem()]
	buffer["DW"] = (*dot)[dot.Elem()]
	buffer["HP"] = 8
	(*ef).Effects = append((*ef).Effects, buffer)
}