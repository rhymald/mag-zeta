package api

import (
	"rhymald/mag-zeta/connect"
	"rhymald/mag-zeta/play"
)

func (loc *Location) GridWriter_ByPush(writeToCache chan map[string][][3]int) {
	for {
		char := <- writeToCache
		connect.WriteChunk((*loc).Writer, char)
		for _, trace := range char {
			x, y := trace[len(trace)-1][1], trace[len(trace)-1][2]
			t, r := play.TAxis(), 700
			connect.ReadRound((*loc).Writer, x, y, r, t)
		}
	}
} 