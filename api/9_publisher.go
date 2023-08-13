package api

import (
	"rhymald/mag-zeta/connect"
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/base"
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
		// whowaswhere := make(map[string][6]int)
		for moment, list := range results {
			for _, each := range list {
				loc.Lock()
				char := (*loc).ByID[each]
				loc.Unlock()
				epoch := base.Epoch()
				even := (epoch/(80*400))%3			
				(*char).Trace.Lock()
				trace, later := (*char).Trace.Erd, (*char).Trace.Snd
				if even == 1 { 
					trace, later = (*char).Trace.Ist, (*char).Trace.Erd 
				} else if even == 2 {
					trace, later = (*char).Trace.Snd, (*char).Trace.Ist
				}			
				(*char).Trace.Unlock()
				current := trace[moment]
				previous := [3]int{}
				iterator := 1
				for { 
					if _, ok := trace[moment-iterator] ; ok { 
						previous = trace[moment-iterator] 
						break
					} else { 
						if moment-iterator < 0 {
							if _, ok := later[moment-iterator+80] ; ok { previous = later[moment-iterator+80] }
						} else if moment-iterator < -80 { break }
						iterator++
					}
				}
				fmt.Println(each, "at", moment, "was in", current, previous)
			}
			break

			// Crushes: concurrent r/w of a map
		}
		// find path cross, if cross collide
		// collide: leftoves of paths after crosses = strength
		// write collision
	}
} 