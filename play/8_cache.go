package play

import (
	"rhymald/mag-zeta/base"
	"sync"
)

type State struct {
	Current *Character
	// Effects map[int]*Effect
	Later struct {
		Time map[string]int
		// Body *base.Stream
		// Energy []*base.Stream
		Life base.Life
		// +Actions 
	}
	Writing struct {
		Time map[string]int
		// Body base.Stream
		// Energy []base.Stream
		Life base.Life
		// +Actions
	}
	sync.Mutex
}

func (c *Character) NewState() *State {
	var buffer State
	c.Lock()
	buffer.Current = c
	buffer.Later.Time = make(map[string]int)
	buffer.Later.Time["life"] = base.EpochNS()
	buffer.Later.Life = *((*c).Life)
	c.Unlock()
	buffer.Writing.Time = make(map[string]int)
	buffer.Writing.Time["life"] = 0 
	buffer.Writing.Life = *(base.MakeLife())
	buffer.Writing.Life.Rate = 0
	return &buffer
}

func (st *State) UpdLife() { // used after write
	(*st).Current.Lock()
	st.Lock()
	timeGape := base.EpochNS() - (*st).Later.Time["life"]
	(*st).Writing.Time["life"] = timeGape
	lifeGape := (*st.Current.Life).Rate - (*st).Later.Life.Rate
	(*st).Writing.Life.Rate = lifeGape
	barriers := make(map[string]int)
	for _, element := range base.ElemList {
		change := (*(*(*st).Current).Life).Barrier[element] - (*st).Writing.Life.Barrier[element]
		if change != 0 { barriers[element] = change }
	}
	(*st).Writing.Life.Barrier = barriers
	(*st).Later.Time["life"] = base.EpochNS()
	(*st).Later.Life = *((*(*st).Current).Life)
	st.Unlock()
	(*st).Current.Unlock()
}
// +write life - tbd in thicket package