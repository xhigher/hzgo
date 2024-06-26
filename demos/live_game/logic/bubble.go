package logic

import (
	"github.com/xhigher/hzgo/demos/live_game/maps"
	"time"
)

type BubbleState int

const (
	BubbleAlive  BubbleState = 0
	BubbleBombed BubbleState = 1
)

type BombResult struct {
	Bombs       []int       `json:"bombs"`
	Boxes       []maps.Site `json:"destroyBoxes"` //需要被摧毁的箱子
	Props       []maps.Site `json:"properties"`   //箱子被炸毁产生的道具
	Areas       []maps.Site `json:"areas"`        //产生爆炸的点
	BubbleCount int         `json:"bombCount"`    //剩余的泡泡数量
}

type Bubble struct {
	Id       int
	site     maps.Site
	color    int
	power    int
	room     *Room
	player   *Player
	bombTime time.Time
	State    BubbleState
}
type BubbleData struct {
	Id          int `json:"i"`
	X           int `json:"x"`
	Y           int `json:"y"`
	Color       int `json:"color"`
	Power       int `json:"power"`
	PlayerIndex int `json:"userIndex"`
}

func (b Bubble) GetData() BubbleData {
	return BubbleData{
		Id:          b.Id,
		X:           b.site.X,
		Y:           b.site.X,
		Color:       b.color,
		Power:       b.power,
		PlayerIndex: b.player.index,
	}
}

func (b *Bubble) CheckBomb(now time.Time) (result *BombResult) {
	if b.State == BubbleAlive && now.After(b.bombTime) {
		result = &BombResult{}
		b.Bomb(result)
	}
	return
}

func (b *Bubble) Bomb(result *BombResult) {
	b.State = BubbleBombed

	//从房间泡泡中删除这个泡泡
	if !b.room.DeleteBubble(b) {
		return
	}
	//炸掉一个泡泡，则要给该用户返回一个泡泡数量
	b.player.bubbleCount++

	//判断自己的位置是否有玩家，有则泡住
	for _, a := range result.Areas {
		if a.Equal(b.site) {
			result.Areas = append(result.Areas, b.site)
		}
	}

	//给爆掉的泡泡数组加入自己本身这个泡泡
	result.Bombs = append(result.Bombs, b.Id)

	//向上爆
	for i := 1; i <= b.power; i++ {
		site := maps.Site{
			X: b.site.X - i,
			Y: b.site.Y,
		}
		if b.room.mapData.IsValidSite(site) {
			if !b.canBomb(site, result) {
				break
			}
		}
	}

	//向下爆
	for i := 1; i <= b.power; i++ {
		site := maps.Site{
			X: b.site.X + i,
			Y: b.site.Y,
		}
		if b.room.mapData.IsValidSite(site) {
			if !b.canBomb(site, result) {
				break
			}
		}
	}

	//向左爆
	for i := 1; i <= b.power; i++ {
		site := maps.Site{
			X: b.site.X,
			Y: b.site.Y + i,
		}
		if b.room.mapData.IsValidSite(site) {
			if !b.canBomb(site, result) {
				break
			}
		}
	}

	//向右爆
	for i := 1; i <= b.power; i++ {
		site := maps.Site{
			X: b.site.X,
			Y: b.site.Y - i,
		}
		if b.room.mapData.IsValidSite(site) {
			if !b.canBomb(site, result) {
				break
			}
		}
	}
}

func (b *Bubble) canBomb(site maps.Site, result *BombResult) bool {
	//判断这个点是否有障碍物，如果有，直接返回false
	if b.room.mapData.ExistObstacle(site) {
		return false
	}
	//判断这个点是否有箱子，如果有，则引爆箱子
	if b.room.mapData.ExistBox(site) {
		for _, bx := range result.Boxes {
			if bx.Equal(site) {
				return false
			}
		}
		result.Boxes = append(result.Boxes, site)
		return false
	}

	//判断该点是否有其他泡泡，如果有，引爆该泡泡
	if yes, b2 := b.room.ExistBubble(site); yes {
		b2.Bomb(result)
		return false
	}

	//如果该爆炸波没有碰到其他物体，则添加到爆炸点
	for _, a := range result.Areas {
		if a.Equal(site) {
			return false
		}
	}
	result.Areas = append(result.Areas, site)

	return true
}
