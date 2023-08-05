package api

import (
	"rhymald/mag-zeta/connect"
)

func (loc *Location) GridWriter_ByPush(writeToCache chan map[string][][3]int) {
	toWrite := make(map[string][][3]int)
	for {
		char := <- writeToCache
		for id, TXYs := range char { 
			if _, ok := toWrite[id]; ok {
				connect.WriteChunk((*loc).Writer, toWrite)
				toWrite = make(map[string][][3]int)
			}
			toWrite[id] = TXYs
			// for _, txy := range TXYs { fmt.Printf("  [GRID] ID: %s => X: %+4d, Y: %+4d @%d\n", id, txy[1], txy[2], txy[0]) }
		}
	}
} 