package maps

import (
	"math/rand"
)

var maps = map[int]Map{}

type Map struct {
	Id int `json:"id"`
	Size Size `json:"size"`
	TileSize Size `json:"tile_size"`
	BornSites []Site `json:"born_sites"`
	Boxes [][]int `json:"boxes"`
	Obstacles [][]int `json:"obstacles"`
}
type Site struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (s Site) Equal(s2 Site) bool{
	return s.X==s2.X && s.Y==s2.Y
}

type Size struct {
	Width int `json:"width"`
	Height int `json:"height"`
}

func GetMap(playerNum int) Map{
	id := (((playerNum-1) / 2)+1) * 200 + rand.Intn(4)
	return maps[id]
}

func addMap(m Map) {
	maps[m.Id] = m
}

func Init(){
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