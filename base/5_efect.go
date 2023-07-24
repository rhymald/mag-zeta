package base

// var Instant_tags = []string{ "HPPP" }
// var Delayed_tags = []string{ "DE", "DW", "DD" }

// var DotTags = []string{ 
// 	"DE", // Dot Element
// 	"DW", // Dot Weight
// 	"DD", // Dot Delay
// }
// var RegenTags = []string{ 
// 	"HPPP",  // HP promille portion
// }

// func GetTagTime(tag string) string {
// 	for _, each := range Instant_tags { if each == tag { return "instant" } }
// 	for _, each := range Delayed_tags { if each == tag { return "delayed" } }
// 	return "ERROR"
// }
// func GetTagType(tag string) string {
// 	for _, each := range DotTags   { if each == tag { return "Self_MakeDot" } }
// 	for _, each := range RegenTags { if each == tag { return "Self_HPRegen" } }
// 	return "ERROR"
// }
// func ValidEffect(effect map[string]int) bool {
// 	first := ""
// 	for tag, _ := range effect {
// 		if len(first) == 0 { first = tag }
// 		if tag != first { return false }
// 	}
// 	return true
// }

// main struct
type Effect struct {
	Time int
	Collision [2]int
	Effects []interface{
		Delayed() int
	}
}
func NewEffect() *Effect {
	buffer := Effect{ Time: Epoch() }
	return &buffer
}
// for _, model := range models {
// 	 if u, ok := model.([]Type1); ok {
// 		 for _, innerUser := range u {
// 			 log.Printf("%#v", innerUser)
// 		 }
// 	 }
// 	 if a, ok := model.([]Type2); ok {
// 	 	 for _, innerArticle := range a {
// 		  	log.Printf("%#v", innerArticle)
// 		 }
// 	 }
// }


// delayed
type Effect_MakeDot struct {
	Dot Dot
	Delay int
}
func (md Effect_MakeDot) Delayed() int { return md.Delay }
func (ef *Effect) Add_Self_MakeDot(dot *Dot) float64 { 
	(*ef).Effects = append((*ef).Effects, 
	Effect_MakeDot{ 
		Dot: *dot,
		Delay: CeilRound(1618/dot.Weight()+1),
	})
	return 1618/dot.Weight()+1
}


// instant
type Effect_HPRegen struct {
	Portion int
}
func (md Effect_HPRegen) Delayed() int { return 0 }
func (ef *Effect) Add_Self_HPRegen(hp int) { (*ef).Effects = append((*ef).Effects, 
	Effect_HPRegen{ 
		Portion: hp,
	})
}


// conditions
// ...
