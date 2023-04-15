package base

import "testing"

func Test_Elements(t *testing.T) {
	t.Logf("Elements:")
	for each, index := range ElemIndex { t.Logf("- %s, rarity: %d", each, index) }
}

func Test_Functions(t *testing.T) {
	for x:=1.0; x*32<2000; x*=1.618 {
		t.Logf("Time:   ms %d: nano %d, waiting %.1f ms", Epoch(), EpochNS(), x*32) ; Wait(x*32)
		t.Logf("Entropy at %d: 1: %.3f, 100: %.3f, 10000: %.3f", Epoch(), Ntrp(1), Ntrp(100), Ntrp(10000))
	}
	for x:=0; x<7; x++ { 
		d:=Near(10*(x+1))
		t.Logf("Random float: %.3f from %d", d, 10*(x+1)) 
		t.Logf(" - Round: %d, Ceil: %d, Flor: %d", Round(d), CeilRound(d), FloorRound(d))
		t.Logf(" - Chaced round: %d, %d, %d, %d, %d, %d, %d, %d, %d, %d", ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d), ChancedRound(d))
	}
}