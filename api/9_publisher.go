package api

import (
	"rhymald/mag-zeta/connect"
	"rhymald/mag-zeta/play"
	"fmt"
)

func (loc *Location) GridWriter_ByPush(writeToCache chan map[string][][3]int) {
	for {
		char := <- writeToCache
		connect.WriteChunk((*loc).Writer, char)
		results := map[int][]string{}
		for _, trace := range char {
			x, y := trace[len(trace)-1][1], trace[len(trace)-1][2]
			t, r := play.TAxis(), 700*2 // replace with step+bodysize
			results = connect.ReadRound((*loc).Writer, x, y, r, t)
			break
		}
		for moment, list := range results { if len(list) < 2 { delete(results, moment) }}
		if len(results) == 0 { continue }
		fmt.Println("Collisions found:", results)

		// find path cross, if cross collide
		// collide: leftoves of paths after crosses = strength
		// write collision
	}
} 