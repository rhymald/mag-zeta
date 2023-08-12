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
			t, r := play.TAxis(), 700*2
			connect.ReadRound((*loc).Writer, x, y, r, t)
			break
		}
		// find path cross, if cross collide
		// collide: leftoves of paths after crosses = strength
		// write collision
	}
} 