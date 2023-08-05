package api

import (
	"rhymald/mag-zeta/connect"
	"rhymald/mag-zeta/base"
	"sync"
)

type Caching struct {
	Cache map[string][][3]int
	sync.Mutex
}

func (loc *Location) GridWriter_ByPush(writeToCache chan map[string][][3]int) {
	interval := 1000000000 / 32
	toWrite := &Caching{ Cache: make(map[string][][3]int) }
	later := base.EpochNS()
	for {
		char := <- writeToCache
		for id, TXYs := range char { 
			if _, ok := (*toWrite).Cache[id]; ok {
				connect.WriteChunk((*loc).Writer, (*toWrite).Cache)
				if base.EpochNS() - later > interval { 
					toWrite.Lock()
					(*toWrite).Cache = make(map[string][][3]int)
					(*toWrite).Cache[id] = TXYs
					toWrite.Unlock()
				} else {
					toWrite.Lock()
					buffer := (*toWrite).Cache[id]
					for _, each := range TXYs { buffer = append(buffer, each) }
					(*toWrite).Cache[id] = buffer
					toWrite.Unlock()
				}
			} else {
				(*toWrite).Cache[id] = TXYs
			}
			// for _, txy := range TXYs { fmt.Printf("  [GRID] ID: %s => X: %+4d, Y: %+4d @%d\n", id, txy[1], txy[2], txy[0]) }
		}
	}
} 