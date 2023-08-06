package play

import (
	"rhymald/mag-zeta/base"
	"sync"
	"math"
)

var timePeriod = 200

type State struct {
	Trace map[int][3]int `json:"Trace"` // time: dir + x + y
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
	buffer.Trace = make(map[int][3]int)
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

func (st *State) Move(writeToCache chan map[string][][3]int) {
	now := base.Epoch()/timePeriod*timePeriod
	st.Lock()
	traceLen := len((*st).Trace)
	if traceLen == 0 { (*st).Trace[now] = [3]int{ base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ) } ; st.Unlock() ; return }
	buffer, latest := (*st).Trace, 0
	for ts, _ := range buffer { if ts > latest { latest = ts } }
	latestStep := (*st).Trace[latest]
	distance := (*st.Current.Atts).Agility // static yet
	angle := float64(latestStep[0])/1000 * math.Pi / 180
 	newstep := [3]int{
		latestStep[0],
		base.Round(float64(latestStep[1]) - 1000*distance*math.Sin(angle)),
		base.Round(float64(latestStep[2]) - 1000*distance*math.Cos(angle)),
	}
	id := (*st).Current.GetID()
	toWrite := make(map[string][][3]int) // id: t, x, y
	for ts := latest+timePeriod ; ts < now ; ts += timePeriod { 
		(*st).Trace[ts] = latestStep 
		toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
	}
	(*st).Trace[now] = newstep
	toWrite[id] = append(toWrite[id], [3]int{now, newstep[1], newstep[2]})
	st.Unlock()
	writeToCache <- toWrite
	base.Wait(math.Phi*1000 / math.Log2(distance+1)) // 1.536 - 0.256
}

func (st *State) Turn(rotate float64, writeToCache chan map[string][][3]int) {
	if math.Abs(rotate) < 1/512 { return }
	now := base.Epoch()/timePeriod*timePeriod
	st.Lock()
	traceLen := len((*st).Trace)
	if traceLen == 0 { (*st).Trace[now] = [3]int{ base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ) } ; st.Unlock() ; return }
	buffer, latest := (*st).Trace, 0
	for ts, _ := range buffer { if ts > latest { latest = ts } }
	latestStep := (*st).Trace[latest]
	distance := (*st.Current.Atts).Agility // static yet
	angle := float64(latestStep[0])/1000 * math.Pi / 180
	newAng := base.Round(angle + rotate*1000)
	for { if newAng > 1000 { newAng += -2000 } else if newAng < -1000 { newAng += 2000 } else { break }}
 	newstep := [3]int{
		newAng,
		latestStep[1],
		latestStep[2],
	}
	id := (*st).Current.GetID()
	toWrite := make(map[string][][3]int) // id: t, x, y
	for ts := latest+timePeriod ; ts < now ; ts += timePeriod { 
		(*st).Trace[ts] = latestStep
		toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
	}
	(*st).Trace[now] = newstep
	toWrite[id] = append(toWrite[id], [3]int{now, newstep[1], newstep[2]})
	st.Unlock()
	writeToCache <- toWrite
	base.Wait(1000/math.Phi/math.Phi / math.Log2(distance+1)) // 0.256 - 0.032
}

func (st *State) Path() [5][2]int {
	period := 4096 // ms
	now := base.Epoch()/timePeriod*timePeriod
	st.Lock() ; trace := (*st).Trace ; st.Unlock()
	if len(trace) == 0 { return [5][2]int{} }
	xs, ys, rs, counter := 0, 0, 0, 0
	xs1, ys1, counter1 := 0, 0, 0
	xs2, ys2, counter2 := 0, 0, 0
	max := 0
	for ts, rXY := range trace {
		if now - ts < period { counter++ ; xs += rXY[1] ; ys += rXY[2] }
		if now - ts < period / 4 { xs2 += rXY[1] ; ys2 += rXY[2] ; counter2++ ; rs += rXY[0] }
		if now - ts < period / 2 { xs1 += rXY[1] ; ys1 += rXY[2] ; counter1++ }
		if ts > max { max = ts }
	}
	latest := trace[max]
	rotate := latest[0] - (rs + latest[0]) / (counter2 + 1)
	for { if rotate > 999 { rotate += -2000 } else if rotate < -1000 { rotate += 2000 } else { break }}
	if counter <= 1 { 
		return [5][2]int{
			[2]int{ latest[0], 0 },
			[2]int{ latest[1], latest[2] },
			[2]int{ latest[1], latest[2] },
			[2]int{ latest[1], latest[2] },
			[2]int{ latest[1], latest[2] },
		}
	}
	if counter1 == 0 {
		return [5][2]int{
			[2]int{ latest[0], rotate },
			[2]int{ latest[1], latest[2] },
			[2]int{ xs/counter, ys/counter },
			[2]int{ xs/counter, ys/counter },
			[2]int{ xs/counter, ys/counter },
		}	
	}
	if counter2 == 0 {
		return [5][2]int{
			[2]int{ latest[0], rotate },
			[2]int{ latest[1], latest[2] },
			[2]int{ xs1/counter1, ys1/counter1 },
			[2]int{ xs1/counter1, ys1/counter1 },
			[2]int{ xs/counter, ys/counter },
		}	
	}
	return [5][2]int{
		[2]int{ latest[0], rotate },
		[2]int{ latest[1], latest[2] },
		[2]int{ xs2/counter2, ys2/counter2 },
		[2]int{ xs1/counter1, ys1/counter1 },
		[2]int{ xs/counter, ys/counter },
	}
}

func (st *State) Collide(object *base.Stream) {
	// TBD
}
