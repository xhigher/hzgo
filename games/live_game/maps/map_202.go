package maps

var map202 = Map{
	Id: 202,
	Size: Size{
		Width: 17,
		Height: 12,
	},
	TileSize: Size{
		Width: 60,
		Height: 60,
	},
	BornSites: []Site{
		{6,5},
		{10,5},
	},
	Boxes: [][]int{
		{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		{0,1,0,1,1,1,1,1,0,1,1,1,1,1,0,1,0},
		{0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,0},
		{0,1,1,1,0,1,1,1,1,1,1,1,0,1,1,1,0},
		{0,1,1,1,1,0,0,1,1,1,0,0,1,1,1,1,0},
		{0,1,1,1,1,0,1,0,0,0,1,0,1,1,1,1,0},
		{0,1,1,1,1,0,0,0,0,0,0,0,1,1,1,1,0},
		{0,1,1,1,0,1,1,1,1,1,1,1,0,1,1,1,0},
		{0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,0},
		{0,1,0,1,1,1,1,1,0,1,1,1,1,1,0,1,0},
		{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
	},
	Obstacles: [][]int{
		{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
		{1,0,1,0,0,0,0,0,1,0,0,0,0,0,1,0,1},
		{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
		{1,0,0,0,1,0,0,0,0,0,0,0,1,0,0,0,1},
		{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
		{1,0,0,0,0,0,0,1,1,1,0,0,0,0,0,0,1},
		{1,0,0,0,0,0,0,1,1,1,0,0,0,0,0,0,1},
		{1,0,0,0,1,0,0,0,0,0,0,0,1,0,0,0,1},
		{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
		{1,0,1,0,0,0,0,0,1,0,0,0,0,0,1,0,1},
		{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
	},
}