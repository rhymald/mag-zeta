package base

import "testing"
import "fmt"

func Test_Elements(t *testing.T) {
	t.Logf("Elements:")
	for each, index := range ElemIndex { t.Logf("- %s, rarity: %d", each, index) }
}

func Test_Functions_Time(t *testing.T) {
	for x:=1.0; x<2000; x*=1.618 {
		t.Logf("Time:   ms %d: nano %d, waiting %.1f ms", 
			Epoch(), EpochNS(), x) 
		Wait(x)
		t.Logf("Entropy at %d: 1: %.3f, 100: %.3f, 10000: %.3f", 
			Epoch(), Ntrp(1), Ntrp(100), Ntrp(10000))
	}
}

func Test_Functions_Rounds(t *testing.T) {
		for x:=0; x<7; x++ { 
		d:=Near(10*(x+1))
		t.Logf("Random float: %.3f from %d", d, 10*(x+1)) 
		t.Logf(" - Round: %d, Ceil: %d, Flor: %d", Round(d), CeilRound(d), FloorRound(d))
		t.Logf(" - Chaced round: %d, %d, %d, %d, %d, %d, %d, %d, %d, %d", 
			ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), 
			ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d))
	}
}

func Test_Streams(t *testing.T) {
	for x:=0; x<9; x++ {
		str := MakeStream()
		_,_, maj, min := str.Constitution(1.1479)
		t.Logf("Stream %d of %s: - Maj:%7.3f | Min:%7.3f | Dev:%7.3f - Len:%7.3f | Mean:%7.3f | Dot:%7.3f \t\t %.1f%% | %.1f%%", 
			x+1, str.Elem(), str.Major(), str.Minor(), str.Deviant(), str.Len(), str.Mean(), str.Dot(), maj*100, min*100)
		dots := " - Dots:"
		for d:=0; d<10; d++ { dd := str.MakeDot() ; Wait(100) ; dots = fmt.Sprintf("%s %+v", dots, *dd) }
		t.Logf("%s", dots)
	}
}