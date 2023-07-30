package play

import (
	"rhymald/mag-zeta/base"
	"sync"
	"math"
)

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
	timeGape := base.EpochNS() - (*st).Later.Time["Life"] 
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
	(*st).Later.Time["Life"] = base.EpochNS()
	(*st).Later.Life = *((*(*st).Current).Life)
	st.Unlock()
	(*st).Current.Unlock()
}
// +write life - tbd in thicket package

func (st *State) Move() {
	st.Lock()
	traceLen := len((*st).Trace)
	if traceLen == 0 { (*st).Trace[base.EpochNS()] = [3]int{ base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ), base.ChancedRound( 2000*base.Rand()-1000 ) } ; st.Unlock() ; return }
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
	(*st).Trace[base.EpochNS()] = newstep
	st.Unlock()
	base.Wait(math.Sqrt2*1000) // 1.4 - 0.25
	// clean from old
	// compare trace, current, db
	// add to trace id current newe
	// add to current if trace later 
}

func (st *State) Turn(angle float64) {
	// turn := base.Round( 1000 / (6 + distance) / 2 )
	// if traceLen % 2 == 1 { turn = -turn } 
	// turn += latestStep[0]
	// for { if turn > 999 { turn += -2000 } else if turn < -999 { turn += 2000 } else { break } }
}

func (st *State) Collide(object base.Stream) {
}