package play

import (
	"rhymald/mag-zeta/base"
	"errors"
	"sync"
	"fmt"
	"crypto/sha512"
  "encoding/binary"
	"github.com/btcsuite/btcutil/base58"
)

type Attributes struct { 
	Vitality float64 `json:"Vitality"`
	Agility float64 `json:"Agility"`
	Capacity float64 `json:"Capacity"`
	Resistance map[string]float64 `json:"Resistance"`
}

type Character struct {
	ID map[string]int `json:"ID"`
	sync.Mutex
	// basics
	Body *base.Stream `json:"Body"`
	Energy []*base.Stream `json:"Energy"`
	// consumables
	Life *base.Life `json:"Life"` 
	Pool map[int]*base.Dot `json:"Pool"`
	// recalculateable stats
	Atts *Attributes `json:"Atts"`
}

func LuckyBorn(time int) int { if time%10 == 0 {return 2} else if time%10 == 9 {return 5} else if time%10 < 5 {return 3} else {return 4} ; return 0}
func (c *Character) IsNPC() bool { return len((*c).Energy) <= 1 }
func (c *Character) GetID() string { 
	in_bytes := make([]byte, 8)
	bid := (*c).ID["Born"]
	aid := (*c).ID["Atts"]
  binary.LittleEndian.PutUint64(in_bytes, uint64(bid))
	str := sha512.Sum512(in_bytes)
	bornID := base58.Encode(str[:])
  binary.LittleEndian.PutUint64(in_bytes, uint64(aid))
	str = sha512.Sum512(in_bytes)
	attsID := base58.Encode(str[:])
	return fmt.Sprintf("%v-%v-%v-%v", bornID[:4], bornID[(len(bornID)-9):len(bornID)], attsID[:1], attsID[(len(attsID)-9):len(attsID)])
	// return bornID
	// return (*c).ID["Born"] 
}

func (c *Character) CalculateAttributes() error {
	c.Lock()
	// if c.GetID() < 1000000 { return errors.New("Character Attributes: Empty character ID.") }
	if len((*c).Energy) == 0 { return errors.New("Character Attributes: Empty character streams.") }
	var buffer Attributes
	buffer.Vitality = (*c).Body.Dot() * 10
	buffer.Agility = (*c).Body.Mean() * 0.7
	buffer.Resistance = make(map[string]float64)
	mod := float64(6 - LuckyBorn((*c).ID["Born"]))
	for _, each := range (*c).Energy { 
		buffer.Resistance[each.Elem()] += each.Mean()
		buffer.Capacity += each.Len() * mod
	}
	(*c).Atts = &buffer
	(*c).ID["Atts"] = base.Epoch()
	c.Unlock()
	return nil
}