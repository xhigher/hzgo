package logic

import "github.com/xhigher/hzgo/games/live_game/maps"

type Bubble struct {
	Id int
	site maps.Site
	color int
	power int
	room *Room
	player *Player
	ct int64
}
type BubbleMsg struct {
	Id int `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Color int `json:"color"`
	Power int `json:"power"`
	Player string `json:"player"`
}

func (b Bubble) GetMsg() BubbleMsg {
	return BubbleMsg{
		Id:    b.Id,
		X:     b.site.X,
		Y:     b.site.X,
		Color: b.color,
		Power: b.power,
		Player:  b.player.id,
	}
}

func (b Bubble) Bomb() {

}


