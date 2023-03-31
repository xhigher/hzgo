package logic

import (
	"github.com/xhigher/hzgo/games/live_game/maps"
	"math/rand"
	"time"
)

type PropType int

const (
	PropBubble PropType = 0
	PropPower PropType = 1
	PropShoes PropType = 2
	PropPin PropType = 3

)

type Prop struct {
	id int
	site maps.Site
	typ PropType
	roomId int
	disappearTime time.Time
}
type PropData struct {
	Id int `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Type PropType `json:"type"`
	Room int `json:"room"`
}

func (p Prop) GetData() PropData {
	return PropData{
		Id:   p.id,
		X:    p.site.X,
		Y:    p.site.Y,
		Type: p.typ,
		Room: p.roomId,
	}
}

type PropWeight struct {
	Name string `json:"name"`
	Type PropType `json:"type"`
	Weight int `json:"weight"`
}

func WeightRandom(items []PropWeight) PropWeight{
	total := 0
	var index []int
	for i, it := range items {
		for j:= 0; j < it.Weight; j++ {
			index = append(index, i)
		}
		total = total + it.Weight
	}
	rn := rand.Intn(total)
	return items[index[rn]]
}