package play

import (
	// "rhymald/mag-zeta/connect"
	"rhymald/mag-zeta/base"
	"sync"
	"math"
	// "fmt"
)

const (
	// 1 minute t-axis / 3: latest = 4 sec
	tAxisStep = 400 //ms for grid
	tRange = 80 //steps per bucket, must be >= Retro
	// tRetro = 3 //steps per retro
	// db supports cleanup every min
)

// var CollideAccelerator = 1600 / tAxisStep

func TAxis() int { return (base.Epoch()/tAxisStep)%tRange }

type Tracing struct {
	Ist map[int][3]int `json:"Ist"` // time % 3 = 1: dir + x + y 
	Snd map[int][3]int `json:"Snd"` // time % 3 = 2: dir + x + y 
	Erd map[int][3]int `json:"Erd"` // time % 3 = 0: dir + x + y 
	sync.Mutex
}

type State struct {
	Trace *Tracing `json:"Trace"`
	Effects map[int]*base.Effect `json:"Effects"`
	Later struct {
		Time map[string]int `json:"Time"`
		// Body *base.Stream
		// Energy []*base.Stream
		Life base.Life `json:"Life"`
		// +Actions 
	} `json:"Later"`
	Writing struct {
		Time map[string]int `json:"Time"`
		// Body base.Stream
		// Energy []base.Stream
		Life base.Life `json:"Life"`
		// +Actions
	} `json:"Writing"`
	sync.Mutex
	Current *Character `json:"Current"`
}

func CleanTrace() *Tracing { return &Tracing{
	Ist: make(map[int][3]int),
	Snd: make(map[int][3]int),
	Erd: make(map[int][3]int),	
}}
func (c *Character) NewState() *State {
	var buffer State
	c.Lock()
	buffer.Current = c
	buffer.Effects = make(map[int]*base.Effect)
	buffer.Later.Time = make(map[string]int)
	buffer.Later.Time["Life"] = base.Epoch()
	buffer.Later.Life = *((*c).Life)
	c.Unlock()
	buffer.Writing.Time = make(map[string]int)
	buffer.Writing.Time["Life"] = 0 
	buffer.Writing.Life = *(base.MakeLife())
	buffer.Writing.Life.Rate = 0
	buffer.Trace = CleanTrace()
	return &buffer
}

// OLD diff rewritten! add to new result!!!
func (st *State) UpdLife() { // used after write
	(*st).Current.Lock()
	st.Lock()
	timeGape := base.Epoch() - (*st).Later.Time["Life"] 
	(*st).Writing.Time["Life"] += timeGape
	lifeGape := (*st.Current.Life).Rate - (*st).Later.Life.Rate
	(*st).Writing.Life.Rate += lifeGape
	// barriers := make(map[string]int)
	for _, element := range base.ElemList {
		change := (*(*(*st).Current).Life).Barrier[element] - (*st).Later.Life.Barrier[element]
		(*st).Writing.Life.Barrier[element] += change 
		if (*st).Writing.Life.Barrier[element] == 0 { delete((*st).Writing.Life.Barrier, element) }
	}
	// (*st).Writing.Life.Barrier = barriers
	(*st).Later.Time["Life"] = base.Epoch()
	(*st).Later.Life = *((*(*st).Current).Life)
	st.Unlock()
	(*st).Current.Unlock()
}
// +write life - tbd in thicket package

func (st *State) Move(rotate float64, step bool, writeToCache chan map[string][][3]int) {
	epoch := base.Epoch()
	now, even := (epoch/tAxisStep)%tRange, (epoch/(tRange*tAxisStep))%3
	(*st).Trace.Lock()
	traceLen := len((*st.Trace).Ist)+len((*st.Trace).Snd)+len((*st.Trace).Erd)
	if traceLen == 0 { 
		if even == 0 {
			(*st.Trace).Erd[now] = [3]int{ 
				base.ChancedRound( 2000*base.Rand()-1000 )/250*250, 
				base.ChancedRound( 2000*base.Rand()-1000 ), 
				base.ChancedRound( 2000*base.Rand()-1000 ),
			}
			(*st).Trace.Unlock()
			return 
		} else if even == 1 {
			(*st.Trace).Ist[now] = [3]int{ 
				base.ChancedRound( 2000*base.Rand()-1000 )/250*250,
				base.ChancedRound( 2000*base.Rand()-1000 ),
				base.ChancedRound( 2000*base.Rand()-1000 ),
			}
			(*st).Trace.Unlock()
			return 
		} else {
			(*st.Trace).Snd[now] = [3]int{ 
				base.ChancedRound( 2000*base.Rand()-1000 )/250*250,
				base.ChancedRound( 2000*base.Rand()-1000 ),
				base.ChancedRound( 2000*base.Rand()-1000 ),
			}
			(*st).Trace.Unlock()
			return 
		}
	}
	later, trace, wipe := (*st.Trace).Snd, (*st.Trace).Erd, &(*st.Trace).Ist
	if even == 1 { 
		later, trace, wipe = (*st.Trace).Erd, (*st.Trace).Ist, &(*st.Trace).Snd 
	} else if even == 2 {
		later, trace, wipe = (*st.Trace).Ist, (*st.Trace).Snd, &(*st.Trace).Erd 
	}
	*wipe = make(map[int][3]int)
	(*st).Trace.Unlock()
	latest, buffer := -tRange, make(map[int][3]int)
	// fmt.Println("Read traces STARTED ================================")
	for ts, each := range later { buffer[ts-tRange] = each }// ; fmt.Println("READing old traces:", even, ts-tRange, each) }
	for ts, each := range trace { buffer[ts] = each }//; fmt.Println("READing current traces:", even, ts, each)}
	// fmt.Println("Read traces ----------------------------")
	// for ts, each := range buffer { fmt.Println("Read traces:", even, ts, each) }
	// fmt.Println("Read traces FINISHED -------------------")
	for ts, _ := range buffer { if ts > latest { latest = ts } }
	latestStep := buffer[latest]
	(*st).Current.Lock()
	distance := (*st.Current.Atts).Agility // static yet
	id := (*st).Current.GetID()
	(*st).Current.Unlock()
	angle := float64(latestStep[0])/1000 // * math.Pi / 180
	newAng := base.Round((angle + rotate)*1000) // * math.Pi / 180
	for { if newAng > 1000 { newAng += -2000 } else if newAng < -1000 { newAng += 2000 } else { break }}
	newstep := [3]int{ newAng, latestStep[1], latestStep[2] }
	if step {
		turn := float64(newAng) / 1000 * math.Pi
		newstep[1] = base.Round(float64(latestStep[1]) + 1000*distance*math.Sin(turn))
		newstep[2] = base.Round(float64(latestStep[2]) + 1000*distance*math.Cos(turn))
		// fmt.Println(angle*180, "--to--", turn/math.Pi*180, "--with--", 1000*math.Sin(turn), 1000*math.Cos(turn))
	}
	toWrite := make(map[string][][3]int) // id: t, x, y
	(*st).Trace.Lock()
	if even == 1 {
		for ts := latest ; ts < now ; ts++ { 
			// if ts >= 0 { (*st).Trace.Ist[ts] = latestStep 
			// 	// fmt.Println("Writing current traces:", even, ts, latestStep) 
			// } else { (*st).Trace.Snd[ts+tRange] = latestStep }//; fmt.Println("Writing old traces:", even, ts+tRange, latestStep)}
			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
		}
		// fmt.Println("Write traces ----------------------------")
		// fmt.Println("WRITTEN trace:", even, now, newstep)
		// fmt.Println("Write traces FINISHED ============================")
		(*st.Trace).Ist[now] = newstep //else { (*st).Trace.Odd[now+(tRange*tAxisStep)/tAxisStep] = newstep }
	} else if even == 2 {
		for ts := latest ; ts < now ; ts++ { 
			// if ts >= 0 { (*st).Trace.Snd[ts] = latestStep 
			// 	// fmt.Println("Writing current traces:", even, ts, latestStep) 
			// } else { (*st).Trace.Erd[ts+tRange] = latestStep }// ; fmt.Println("Writing old traces:", even, ts+tRange, latestStep)}
			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
		}
		// fmt.Println("Write traces ----------------------------")
		// fmt.Println("WRITTEN trace:", even, now, newstep)
		// fmt.Println("Write traces FINISHED ============================")
		(*st.Trace).Snd[now] = newstep
	} else {
		for ts := latest ; ts < now ; ts++ { 
			// if ts >= 0 { (*st).Trace.Erd[ts] = latestStep 
			// 	// fmt.Println("Writing current traces:", even, ts, latestStep) 
			// } else { (*st).Trace.Ist[ts+tRange] = latestStep }// ; fmt.Println("Writing old traces:", even, ts+tRange, latestStep)}
			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
		}
		// fmt.Println("Write traces ----------------------------")
		// fmt.Println("WRITTEN trace:", even, now, newstep)
		// fmt.Println("Write traces FINISHED ============================")
		(*st.Trace).Erd[now] = newstep
	}
	(*st).Trace.Unlock()
	toWrite[id] = append(toWrite[id], [3]int{now, newstep[1], newstep[2]})
	writeToCache <- toWrite
	base.Wait(float64(tAxisStep)*math.Pi)// / math.Log2(distance+1)) // 1.536 - 0.256
	// connect.ReadRound(world. )
}

// func (st *State) Turn(rotate float64, writeToCache chan map[string][][3]int) {
// 	// if math.Abs(rotate) < 1/512 { return }
// 	epoch := base.Epoch()
// 	now, even := (epoch%(tRange*tAxisStep))/tAxisStep, (epoch/(tRange*tAxisStep))%2 == 0
// 	(*st).Trace.Lock()
// 	traceLen := len((*st).Trace.Odd)+len((*st).Trace.Even)
// 	if traceLen == 0 { 
// 		if even {
// 			(*st).Trace.Even[now] = [3]int{ 
// 				base.ChancedRound( 2000*base.Rand()-1000 ), 
// 				base.ChancedRound( 2000*base.Rand()-1000 ), 
// 				base.ChancedRound( 2000*base.Rand()-1000 ),
// 			}
// 			(*st).Trace.Unlock()
// 			return 
// 		} else {
// 			(*st).Trace.Odd[now] = [3]int{ 
// 				base.ChancedRound( 2000*base.Rand()-1000 ),
// 				base.ChancedRound( 2000*base.Rand()-1000 ),
// 				base.ChancedRound( 2000*base.Rand()-1000 ),
// 			}
// 			(*st).Trace.Unlock()
// 			return 
// 		}
// 	}
// 	trace, later, latest, buffer := (*st).Trace.Odd, (*st).Trace.Even, -(tRange*tAxisStep)*2/tAxisStep-1, map[int][3]int{}
// 	if even { trace, later = (*st).Trace.Even, (*st).Trace.Odd }
// 	for ts, each := range later { buffer[ts-(tRange*tAxisStep)/tAxisStep] = each }
// 	for ts, each := range trace { if ts < now { buffer[ts] = each } else { 
// 		buffer[ts-(tRange*tAxisStep)/tAxisStep] = each 
// 		if even { delete((*st).Trace.Even, ts) } else { delete((*st).Trace.Odd, ts) }
// 	}}
// 	for ts, _ := range buffer { if ts > latest { latest = ts } }
// 	latestStep := buffer[latest]
// 	(*st).Current.Lock()
// 	// distance := (*st.Current.Atts).Agility // to be replaced 
// 	id := (*st).Current.GetID()
// 	(*st).Current.Unlock()
// 	angle := float64(latestStep[0])/1000 // * math.Pi / 180
// 	newAng := base.Round((angle + rotate)*1000) // * math.Pi / 180
// 	// for { if newAng > 1000 { newAng += -2000 } else if newAng < -1000 { newAng += 2000 } else { break }}
//  newstep := [3]int{
// 		newAng,
// 		latestStep[1],
// 		latestStep[2],
// 	}
// 	toWrite := make(map[string][][3]int) // id: t, x, y
// 	if even {
// 		for ts := latest ; ts < now ; ts++ { 
// 			if ts > 0 { (*st).Trace.Even[ts] = latestStep } else { (*st).Trace.Odd[ts+(tRange*tAxisStep)/tAxisStep] = latestStep }
// 			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
// 		}
// 		if now > 0 { (*st).Trace.Even[now] = newstep } else { (*st).Trace.Odd[now+(tRange*tAxisStep)/tAxisStep] = newstep }
// 		// (*st).Trace[now] = newstep
// 	} else {
// 		for ts := latest ; ts < now ; ts++ { 
// 			if ts > 0 { (*st).Trace.Odd[ts] = latestStep } else { (*st).Trace.Even[ts+(tRange*tAxisStep)/tAxisStep] = latestStep }
// 			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
// 		}
// 		if now > 0 { (*st).Trace.Odd[now] = newstep } else { (*st).Trace.Even[now+(tRange*tAxisStep)/tAxisStep] = newstep }
// 	}
// 	toWrite[id] = append(toWrite[id], [3]int{now, newstep[1], newstep[2]})
// 	// fmt.Println((*st).Trace)
// 	(*st).Trace.Unlock()
// 	writeToCache <- toWrite
// 	base.Wait(float64(tAxisStep))// / math.Log2(distance+1)) // 0.256 - 0.032
// }

func (st *State) Path() [5][2]int {
	// period := tRetro // steps
	epoch := base.Epoch()
	even := (epoch/(tRange*tAxisStep))%3
	(*st).Trace.Lock()
	trace, later := (*st.Trace).Erd, (*st.Trace).Snd
	if even == 1 { 
		trace, later = (*st.Trace).Ist, (*st.Trace).Erd 
	} else if even == 2 {
		trace, later = (*st.Trace).Snd, (*st.Trace).Ist
	}
	(*st).Trace.Unlock()
	if len(trace)+len(later) == 0 { return [5][2]int{} }
	buffer := make(map[int][3]int)
	for ts, each := range later { buffer[ts-tRange] = each }
	for ts, each := range trace { buffer[ts] = each }
	tIndex, idx := []int{}, 0
	for i:=tRange ; i>-tRange ; i-- {
		if _, ok := buffer[i] ; ok { tIndex = append(tIndex, i) ; idx++ }
		if idx >= 4 { break }
	}
	// xs1, ys1, counter1 := 0, 0, 0
	// xs2, ys2, counter2 := 0, 0, 0
	// xs3, ys3, counter3 := 0, 0, 0
	// max := -tRange ; for ts, _ := range buffer { if ts > max { max = ts } }
	// latest := buffer[max]
	// delete(buffer, max)
	// for ts, rXY := range buffer {
	// 	if ((now-ts)*3+1) / period == 0 { 
	// 		xs1 += rXY[1] ; ys1 += rXY[2] ; counter1++ ; rs += rXY[0]; counter++
	// 	} else if ((now-ts)*3+1) / period == 1 { 
	// 		xs2 += rXY[1] ; ys2 += rXY[2] ; counter2++ ; rs += rXY[0]; counter++
	// 	} else if ((now-ts)*3+1) / period == 2 { 
	// 		xs3 += rXY[1] ; ys3 += rXY[2] ; counter3++ ; rs += rXY[0]; counter++
	// 	}
	// }
	// rotate := rs - latest[0]*counter
// 	for { if rotate > 999 { rotate += -2000 } else if rotate < -1000 { rotate += 2000 } else { break }}
	if idx <= 1 { 
		return [5][2]int{
			[2]int{ buffer[tIndex[0]][0], 0 },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
		}
	}
	if idx == 2 {
		return [5][2]int{
			[2]int{ buffer[tIndex[0]][0], -buffer[tIndex[0]][0]+buffer[tIndex[1]][0] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
			[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
			[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
		}	
	}
	if idx == 3 {
		return [5][2]int{
			[2]int{ buffer[tIndex[0]][0], -buffer[tIndex[0]][0]+buffer[tIndex[1]][0] },
			[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
			[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
			[2]int{ buffer[tIndex[2]][1], buffer[tIndex[2]][2] },
			[2]int{ buffer[tIndex[2]][1], buffer[tIndex[2]][2] },
		}	
	}
	return [5][2]int{
		[2]int{ buffer[tIndex[0]][0], -buffer[tIndex[0]][0]+buffer[tIndex[1]][0] },
		[2]int{ buffer[tIndex[0]][1], buffer[tIndex[0]][2] },
		[2]int{ buffer[tIndex[1]][1], buffer[tIndex[1]][2] },
		[2]int{ buffer[tIndex[2]][1], buffer[tIndex[2]][2] },
		[2]int{ buffer[tIndex[3]][1], buffer[tIndex[3]][2] },
	}
}

func (st *State) Collide(object *base.Stream) {
	// TBD
}
