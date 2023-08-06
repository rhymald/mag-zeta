package play

import (
	"rhymald/mag-zeta/base"
	"sync"
	"math"
)

const (
	tAxisStep = 400 //ms for grid, keep it <500
	tRange = 16000 //ms per bucket, must be >= x2 Retro
	tRetro = 4096 //ms let it be %4
)

type State struct {
	Trace struct {
		Odd map[int][3]int `json:"Odd"`  // time !% 2: dir + x + y 
		Even map[int][3]int `json:"Even"` // time % 2: dir + x + y 
	} `json:"Trace"`
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
	buffer.Trace.Odd = make(map[int][3]int)
	buffer.Trace.Even = make(map[int][3]int)
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
	epoch := base.Epoch()
	now, even := (epoch%tRange)/tAxisStep, (epoch/tRange)%2 == 0
	st.Lock()
	traceLen := len((*st).Trace.Odd)+len((*st).Trace.Even)
	if traceLen == 0 { 
		if even {
			(*st).Trace.Even[now] = [3]int{ base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ) } ; st.Unlock() ; return 
		} else {
			(*st).Trace.Odd[now]  = [3]int{ base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ) } ; st.Unlock() ; return 
		}
	}
	trace, later, latest, buffer := (*st).Trace.Odd, (*st).Trace.Even, -tRange*2/tAxisStep-1, map[int][3]int{}
	if even { trace, later = (*st).Trace.Even, (*st).Trace.Odd }
	for ts, each := range later { buffer[ts-tRange/tAxisStep] = each }
	for ts, each := range trace { if ts < now { buffer[ts] = each } else { 
		buffer[ts-tRange/tAxisStep] = each 
		if even { delete((*st).Trace.Even, ts) } else { delete((*st).Trace.Odd, ts) }
	}}
	for ts, _ := range buffer { if ts > latest { latest = ts } }
	latestStep := buffer[latest]
	distance := (*st.Current.Atts).Agility // static yet
	angle := float64(latestStep[0])/1000 * math.Pi / 180
	newstep := [3]int{
		latestStep[0],
		base.Round(float64(latestStep[1]) - 1000*distance*math.Sin(angle)),
		base.Round(float64(latestStep[2]) - 1000*distance*math.Cos(angle)),
 	}
	id := (*st).Current.GetID()
	toWrite := make(map[string][][3]int) // id: t, x, y
	if even {
		for ts := latest ; ts < now ; ts++ { 
			if ts > 0 { (*st).Trace.Even[ts] = latestStep } else { (*st).Trace.Odd[ts+tRange/tAxisStep] = latestStep }
			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
		}
		if now > 0 { (*st).Trace.Even[now] = latestStep } else { (*st).Trace.Odd[now+tRange/tAxisStep] = latestStep }
		// (*st).Trace[now] = newstep
	} else {
		for ts := latest ; ts < now ; ts++ { 
			if ts > 0 { (*st).Trace.Odd[ts] = latestStep } else { (*st).Trace.Even[ts+tRange/tAxisStep] = latestStep }
			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
		}
		if now > 0 { (*st).Trace.Odd[now] = latestStep } else { (*st).Trace.Even[now+tRange/tAxisStep] = latestStep }
	}
	toWrite[id] = append(toWrite[id], [3]int{now, newstep[1], newstep[2]})
	st.Unlock()
	writeToCache <- toWrite
	base.Wait(math.Phi*1000 / math.Log2(distance+1)) // 1.536 - 0.256
}

func (st *State) Turn(rotate float64, writeToCache chan map[string][][3]int) {
	if math.Abs(rotate) < 1/512 { return }
	epoch := base.Epoch()
	now, even := (epoch%tRange)/tAxisStep, (epoch/tRange)%2 == 0
	st.Lock()
	traceLen := len((*st).Trace.Odd)+len((*st).Trace.Even)
	if traceLen == 0 { 
		if even {
			(*st).Trace.Even[now] = [3]int{ base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ) } ; st.Unlock() ; return 
		} else {
			(*st).Trace.Odd[now]  = [3]int{ base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ) } ; st.Unlock() ; return 
		}
	}
	trace, later, latest, buffer := (*st).Trace.Odd, (*st).Trace.Even, -tRange*2/tAxisStep-1, map[int][3]int{}
	if even { trace, later = (*st).Trace.Even, (*st).Trace.Odd }
	for ts, each := range later { buffer[ts-tRange/tAxisStep] = each }
	for ts, each := range trace { if ts < now { buffer[ts] = each } else { 
		buffer[ts-tRange/tAxisStep] = each 
		if even { delete((*st).Trace.Even, ts) } else { delete((*st).Trace.Odd, ts) }
	}}
	for ts, _ := range buffer { if ts > latest { latest = ts } }
	latestStep := buffer[latest]
	distance := (*st.Current.Atts).Agility // to be replaced 
	angle := float64(latestStep[0])/1000 // * math.Pi / 180
	newAng := base.Round(angle + rotate)*1000 // * math.Pi / 180
	for { if newAng > 1000 { newAng += -2000 } else if newAng < -1000 { newAng += 2000 } else { break }}
 	newstep := [3]int{
		newAng,
		latestStep[1],
		latestStep[2],
	}
	id := (*st).Current.GetID()
	toWrite := make(map[string][][3]int) // id: t, x, y
	if even {
		for ts := latest ; ts < now ; ts++ { 
			if ts > 0 { (*st).Trace.Even[ts] = latestStep } else { (*st).Trace.Odd[ts+tRange/tAxisStep] = latestStep }
			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
		}
		if now > 0 { (*st).Trace.Even[now] = latestStep } else { (*st).Trace.Odd[now+tRange/tAxisStep] = latestStep }
		// (*st).Trace[now] = newstep
	} else {
		for ts := latest ; ts < now ; ts++ { 
			if ts > 0 { (*st).Trace.Odd[ts] = latestStep } else { (*st).Trace.Even[ts+tRange/tAxisStep] = latestStep }
			toWrite[id] = append(toWrite[id], [3]int{ts, latestStep[1], latestStep[2]})
		}
		if now > 0 { (*st).Trace.Odd[now] = latestStep } else { (*st).Trace.Even[now+tRange/tAxisStep] = latestStep }
	}
	toWrite[id] = append(toWrite[id], [3]int{now, newstep[1], newstep[2]})
	st.Unlock()
	writeToCache <- toWrite
	base.Wait(1000/math.Phi/math.Phi / math.Log2(distance+1)) // 0.256 - 0.032
}

// Need to fix! It does not work for /around/...
func (st *State) Path() [5][2]int {
	period := tRetro // ms
	epoch := base.Epoch()
	now, even := (epoch%tRange)/tAxisStep, (epoch/tRange)%2 == 0
	st.Lock()
	// trace := (*st).Trace
	trace, later, buffer := (*st).Trace.Odd, (*st).Trace.Even, map[int][3]int{}
	if even { trace, later = (*st).Trace.Even, (*st).Trace.Odd }
	if len(trace)+len(later) == 0 { st.Unlock() ; return [5][2]int{} }
	for ts, each := range later { buffer[ts-tRange/tAxisStep] = each }
	for ts, each := range trace { if ts < now { buffer[ts] = each } else { 
		buffer[ts-tRange/tAxisStep] = each 
		if even { delete((*st).Trace.Even, ts) } else { delete((*st).Trace.Odd, ts) }
	}}
	// for ts, _ := range trace { if ts > latest { latest = ts } }
	// latestStep := trace[latest]
	st.Unlock()
	xs, ys, rs, counter := 0, 0, 0, 0
	xs1, ys1, counter1 := 0, 0, 0
	xs2, ys2, counter2 := 0, 0, 0
	max := -tRange*2/tAxisStep - 1
	for ts, rXY := range buffer {
		if ts > max { max = ts }
		if now - ts < period { counter++ ; xs += rXY[1] ; ys += rXY[2] }
		if now - ts < period / 4 { xs2 += rXY[1] ; ys2 += rXY[2] ; counter2++ ; rs += rXY[0] }
		if now - ts < period / 2 { xs1 += rXY[1] ; ys1 += rXY[2] ; counter1++ }
	}
	latest := buffer[max]
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
