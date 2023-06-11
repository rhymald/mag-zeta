package base

import (
	"fmt"
)

type Dot map[string]int 

// NEW
func (str *Stream) MakeDot() *Dot { return &Dot{ str.Elem(): FloorRound(Ntrp( str.Dot()*1000 )) }}


// READ
func (d *Dot) ToStr() string { return fmt.Sprintf("%s|%d", d.Elem(), (*d)[d.Elem()] ) }
func (d *Dot) Weight() float64 { return float64((*d)[d.Elem()]) / 1000 }
func (d *Dot) Elem() string { 
	if len(*d) != 1 { return "Error" }
	for elem, _ := range *d { return elem } 
	return "Error"
}