package maps

import (
	"math/rand"
)

var maps = map[int]MapData{}

type MapData struct {
	Id        int     `json:"id"`
	Size      Size    `json:"mapSize"`
	TileSize  Size    `json:"tileSize"`
	BornSites []Site  `json:"userStartPos"`
	Boxes     [][]int `json:"boxes"`
	Obstacles [][]int `json:"obstacles"`
}
type Site struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (s Site) Equal(s2 Site) bool {
	return s.X == s2.X && s.Y == s2.Y
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func GetMap(playerNum int) MapData {
	id := (((playerNum-1)/2)+1)*200 + rand.Intn(4)
	return maps[id]
}

func addMap(m MapData) {
	maps[m.Id] = m
}

func Init() {
	addMap(map200)
	addMap(map201)
	addMap(map202)
	addMap(map203)

	addMap(map400)
	addMap(map401)
	addMap(map402)
	addMap(map403)

	addMap(map600)
	addMap(map601)
	addMap(map602)
	addMap(map603)
}

func (m MapData) GetValidSites(origins []Site) (results []Site) {
	for _, s := range origins {
		if s.X >= 0 && s.Y >= 0 && s.X < m.Size.Height && s.Y < m.Size.Width {
			results = append(results, s)
		}
	}
	return
}

func (m MapData) IsValidSite(s Site) bool {
	if s.X >= 0 && s.Y >= 0 && s.X < m.Size.Height && s.Y < m.Size.Width {
		return true
	}
	return false
}

func (m MapData) ExistObstacle(s Site) bool {
	if m.Obstacles[s.Y] != nil && m.Obstacles[s.Y][s.X] == 1 {
		return true
	}
	return false
}

func (m MapData) ExistBox(s Site) bool {
	if m.Boxes[s.Y] != nil && m.Boxes[s.Y][s.X] == 1 {
		return true
	}
	return false
}

func (m MapData) IsEmptySite(s Site) bool {
	if !m.ExistObstacle(s) && !m.ExistBox(s) {
		return true
	}
	return false
}
