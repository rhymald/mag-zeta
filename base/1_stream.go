package base

import "math"

type Stream map[string][3]int

// NEW
func MakeStream() *Stream {
	maj, min, dev := 0.0, 0.0, 0.0
	for x:=0; x<4; x++ { maj += Rand() ; min += Rand() ; dev += Rand() }
	leng := Vector(maj, min, dev) 
	return &Stream{ ElemList[0]: [3]int{ CeilRound( maj/leng * 1000 ), CeilRound( min/leng * 1000 ), CeilRound( dev/leng * 1000 ) }}
}


// MOD
func (str *Stream) Len() float64 { return Vector( str.Major(), str.Minor(), str.Deviant()) }
func (str *Stream) Mean() float64 { return math.Pi/( 1/str.Major() + 1/str.Minor() + 1/str.Deviant() ) }
func (str *Stream) Dot() float64 { a := math.Log2(str.Mean()+2)/math.Log2(7) ; return math.Pow(a,a) }

func (str *Stream) Major() float64 { return float64((*str)[str.Elem()][0]) / 1000 }
func (str *Stream) Minor() float64 { return float64((*str)[str.Elem()][1]) / 1000 }
func (str *Stream) Deviant() float64 { return float64((*str)[str.Elem()][2]) / 1000 }

func (str *Stream) Elem() string { 
	if len(*str) != 1 { return "Error" }
	for elem, _ := range *str { return elem } 
	return "Error"
}

func (str *Stream) Des() float64 {
	if str.Elem() == ElemList[0] || str.Elem() == ElemList[2] { return str.Major() }
	if str.Elem() == ElemList[1] { return str.Minor() }
	return 0
}

func (str *Stream) Alt() float64 {
	if str.Elem() == ElemList[0] { return str.Deviant() }
	if str.Elem() == ElemList[2] { return str.Minor() }
	if str.Elem() == ElemList[1] { return str.Major() }
	return 0
}

func (str *Stream) Cre() float64 {
	if str.Elem() == ElemList[0] { return str.Minor() }
	if str.Elem() == ElemList[2] || str.Elem() == ElemList[1] { return str.Deviant() }
	return 0
}

func (str *Stream) Constitution(mod float64) (int, int, float64, float64) {
	a, b, c, beeg, smol, tiny := str.Major(), str.Minor(), str.Deviant(), 0.0, 0.0, 0.0
	maj, min := -1, -1
	if a > b && a > c { maj = 0 ; if b < c { min = 2 ; smol, tiny = c, b } else { min = 1 ; smol, tiny = b, c }} ; beeg = a 
	if b > a && b > c { maj = 1 ; if a < c { min = 2 ; smol, tiny = c, a } else { min = 0 ; smol, tiny = c, a }} ; beeg = b
	if c > b && c > a { maj = 2 ; if b < a { min = 0 ; smol, tiny = a, b } else { min = 1 ; smol, tiny = b, a }} ; beeg = c
	return maj, min, math.Log2(beeg/smol)/math.Log2(mod), math.Log2(smol/tiny)/math.Log2(mod)
}