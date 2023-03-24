package events

//{"e":"loginSuccess","_id":"63f08b2d8b2f2a4e2b30ddd9","x":0,"y":0,"role":400,"sex":0,"group":0,"nickname":"xhigher","headImgUrl":"https://thirdwx.qlogo.cn/mmopen/vi_32/KMl1PYjbMoDfeslGzgLfHbxicEuRIJfvwr82SooU7xmW0LX4VEFBRj7wFHADEHcP1GYX9WB620XxZpibhzDP3P3g/132","status":0,"moveInterval":300,"bombCount":2,"thumbtackCount":0,"isReady":false,"level":0}
type LoginSuccess struct {
	E              string `json:"e"`
	Id             string `json:"_id"`
	X              int    `json:"x"`
	Y              int    `json:"y"`
	Role           int    `json:"role"`
	Sex            int    `json:"sex"`
	Group          int    `json:"group"`
	Nickname       string `json:"nickname"`
	HeadImgUrl     string `json:"headImgUrl"`
	Status         int    `json:"status"`
	MoveInterval   int    `json:"moveInterval"`
	BombCount      int    `json:"bombCount"`
	ThumbtackCount int    `json:"thumbtackCount"`
	IsReady        bool   `json:"isReady"`
	Level          int    `json:"level"`
}

//{"e":"matching","self":{"_id":"63f08b2d8b2f2a4e2b30ddd9","x":0,"y":0,"role":300,"sex":0,"group":0,"nickname":"xhigher","headImgUrl":"https://thirdwx.qlogo.cn/mmopen/vi_32/KMl1PYjbMoDfeslGzgLfHbxicEuRIJfvwr82SooU7xmW0LX4VEFBRj7wFHADEHcP1GYX9WB620XxZpibhzDP3P3g/132","status":0,"moveInterval":300,"bombCount":2,"thumbtackCount":0,"isReady":false,"level":0},"users":[{"_id":"63ec3b6a97611132f5ef419c","x":0,"y":0,"role":400,"sex":0,"group":0,"nickname":"机器人200","headImgUrl":"https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKVEyX2hqUPnshYEiarvhh1FtybiapVsBf4SY8ibJy8X6ial9LXUYkfLY0w5JicHHyOAFZUMS8g3zicibDvA/132","status":0,"moveInterval":300,"bombCount":2,"thumbtackCount":0,"isReady":false,"level":0},{"_id":"63ec3b6a97611132f5ef41af","x":0,"y":0,"role":300,"sex":1,"group":0,"nickname":"机器人201","headImgUrl":"https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKVEyX2hqUPnshYEiarvhh1FtybiapVsBf4SY8ibJy8X6ial9LXUYkfLY0w5JicHHyOAFZUMS8g3zicibDvA/132","status":0,"moveInterval":300,"bombCount":2,"thumbtackCount":0,"isReady":false,"level":0},{"_id":"63ec3b6a97611132f5ef41c2","x":0,"y":0,"role":400,"sex":0,"group":0,"nickname":"机器人202","headImgUrl":"https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKVEyX2hqUPnshYEiarvhh1FtybiapVsBf4SY8ibJy8X6ial9LXUYkfLY0w5JicHHyOAFZUMS8g3zicibDvA/132","status":0,"moveInterval":300,"bombCount":2,"thumbtackCount":0,"isReady":false,"level":0},{"_id":"63ec3b6a97611132f5ef41d5","x":0,"y":0,"role":300,"sex":1,"group":0,"nickname":"机器人203","headImgUrl":"https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKVEyX2hqUPnshYEiarvhh1FtybiapVsBf4SY8ibJy8X6ial9LXUYkfLY0w5JicHHyOAFZUMS8g3zicibDvA/132","status":0,"moveInterval":300,"bombCount":2,"thumbtackCount":0,"isReady":false,"level":0},{"_id":"63ec3b6a97611132f5ef41e8","x":0,"y":0,"role":400,"sex":0,"group":0,"nickname":"机器人204","headImgUrl":"https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKVEyX2hqUPnshYEiarvhh1FtybiapVsBf4SY8ibJy8X6ial9LXUYkfLY0w5JicHHyOAFZUMS8g3zicibDvA/132","status":0,"moveInterval":300,"bombCount":2,"thumbtackCount":0,"isReady":false,"level":0}]}
type Matching struct {
	E    string `json:"e"`
	Self interface{} `json:"self"`
	Users []interface{} `json:"users"`
}

//{"e":"gameStart","data":{"roomId":1,"type":0,"roomType":0,"status":0,"users":[{"_id":"63f08b2d8b2f2a4e2b30ddd9","x":8,"y":6,"role":300,"sex":0,"group":0,"nickname":"xhigher","headImgUrl":"https://thirdwx.qlogo.cn/mmopen/vi_32/KMl1PYjbMoDfeslGzgLfHbxicEuRIJfvwr82SooU7xmW0LX4VEFBRj7wFHADEHcP1GYX9WB620XxZpibhzDP3P3g/132","status":0,"moveInterval":300,"bombCount":2,"thumbtackCount":0,"isReady":false,"level":0,"roomId":1},{"_id":"63ec3b6a97611132f5ef439d","x":8,"y":10,"role":300,"sex":1,"group":1,"nickname":"机器人227","headImgUrl":"https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKVEyX2hqUPnshYEiarvhh1FtybiapVsBf4SY8ibJy8X6ial9LXUYkfLY0w5JicHHyOAFZUMS8g3zicibDvA/132","status":0,"moveInterval":300,"bombCount":1,"thumbtackCount":1,"isReady":false,"level":0,"roomId":1},{"_id":"63ec3b6a97611132f5ef43b0","x":12,"y":6,"role":400,"sex":0,"group":2,"nickname":"机器人228","headImgUrl":"https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKVEyX2hqUPnshYEiarvhh1FtybiapVsBf4SY8ibJy8X6ial9LXUYkfLY0w5JicHHyOAFZUMS8g3zicibDvA/132","status":0,"moveInterval":300,"bombCount":1,"thumbtackCount":1,"isReady":false,"level":0,"roomId":1}],"map":{"id":403,"mapSize":{"width":21,"height":16},"tileSize":{"width":60,"height":60},"userStartPos":[[8,6],[8,10],[12,6],[12,10]],"boxes":[[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,0,0],[0,1,1,1,1,1,1,0,1,1,1,1,1,0,1,1,1,1,1,1,0],[0,1,1,0,1,1,1,1,1,1,1,1,1,1,1,1,1,0,1,1,0],[0,1,1,1,1,1,0,0,0,1,0,1,0,0,0,1,1,1,1,1,0],[0,1,1,1,1,1,0,1,0,1,1,1,0,1,0,1,1,1,1,1,0],[0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,0],[0,1,0,0,1,1,0,1,1,0,0,0,1,1,0,1,1,0,0,1,0],[0,1,1,1,1,1,1,1,1,0,0,0,1,1,1,1,1,1,1,1,0],[0,1,1,1,1,1,0,1,0,1,1,1,0,1,0,1,1,1,1,1,0],[0,1,1,1,1,1,0,0,0,1,0,1,0,0,0,1,1,1,1,1,0],[0,1,1,0,1,1,1,1,1,1,1,1,1,1,1,1,1,0,1,1,0],[0,1,1,1,1,1,1,0,1,1,1,1,1,0,1,1,1,1,1,1,0],[0,0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]]},"bombs":[],"properties":[]}}
type GameStart struct {
	E    string `json:"e"`
	Data struct {
		RoomId   int `json:"roomId"`
		Type     int `json:"type"`
		RoomType int `json:"roomType"`
		Status   int `json:"status"`
		Users    []interface{} `json:"users"`
		Map interface{} `json:"map"`
		Bombs      []interface{} `json:"bombs"`
		Properties []interface{} `json:"properties"`
	} `json:"data"`
}


//{"e":"loadProcess","i":1,"process":100}
type LoadProcess struct {
	E       string `json:"e"`
	I       int    `json:"i"`
	Process int    `json:"process"`
}

//{"e":"gameReady"}
type GameReady struct {
	E string `json:"e"`
}

//{"e":"moveStop","i":2}
type MoveStop struct {
	E string `json:"e"`
	I int    `json:"i"`
}

//{"e":"move","i":2,"x":12,"y":5,"t":300}
type Move struct {
	E string `json:"e"`
	I int    `json:"i"`
	X int    `json:"x"`
	Y int    `json:"y"`
	T int    `json:"t"`
}
//{"e":"createBomb","userIndex":2,"i":0,"x":12,"y":5,"color":0,"power":1}
//{"e":"createBomb","userIndex":2,"i":101,"x":2,"y":7,"color":0,"power":3}
//{"e":"createBomb","userIndex":1,"i":105,"x":1,"y":10,"color":0,"power":4}
type CreateBomb struct {
	E         string `json:"e"`
	UserIndex int    `json:"userIndex"`
	I         int    `json:"i"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Color     int    `json:"color"`
	Power     int    `json:"power"`
}

//{"e":"bomb","bombCount":2,"bombs":[0],"destroyBoxes":[[11,5],[12,4]],"properties":[]}
//{"e":"bomb","bombCount":2,"bombs":[1],"destroyBoxes":[[9,11],[8,12]],"properties":[{"i":0,"x":8,"y":12,"type":0}]}
//{"e":"bomb","bombCount":2,"bombs":[2],"destroyBoxes":[[15,5],[14,4]],"properties":[{"i":1,"x":15,"y":5,"type":2},{"i":2,"x":14,"y":4,"type":1}]}
//{"e":"bomb","bombCount":2,"bombs":[6],"destroyBoxes":[[7,12],[7,10]],"properties":[{"i":7,"x":7,"y":12,"type":1}]}
//{"e":"bomb","bombCount":2,"bombs":[7],"destroyBoxes":[[4,11],[5,12],[5,10]],"properties":[{"i":8,"x":4,"y":11,"type":0},{"i":9,"x":5,"y":12,"type":1},{"i":10,"x":5,"y":10,"type":0}]}
//{"e":"bomb","bombCount":2,"bombs":[62,63,64,66,65,67,68],"destroyBoxes":[[2,11],[1,12],[8,9],[1,8]],"properties":[]}
//{"e":"bomb","bombCount":2,"bombs":[69,71,74,73,75,72,70],"destroyBoxes":[[2,7],[2,6],[2,5],[4,2]],"properties":[{"i":39,"x":2,"y":7,"type":1},{"i":40,"x":2,"y":6,"type":2},{"i":41,"x":2,"y":5,"type":2},{"i":42,"x":4,"y":2,"type":0}]}
type Bomb struct {
	E            string        `json:"e"`
	BombCount    int           `json:"bombCount"`
	Bombs        []int         `json:"bombs"`
	DestroyBoxes [][]int       `json:"destroyBoxes"`
	Properties   []interface{} `json:"properties"`
}

//{"e":"disappearProperty","i":0,"isEat":true,"userIndex":1}
//{"e":"disappearProperty","i":4,"isEat":true,"userIndex":2}
//{"e":"disappearProperty","i":1}
type DisappearProperty struct {
	E         string `json:"e"`
	I         int    `json:"i"`
	IsEat     bool   `json:"isEat"`
	UserIndex int    `json:"userIndex"`
}

//{"e":"changeUserStatus","index":0,"x":8,"y":6,"status":-1,"properties":[]}
//{"e":"changeUserStatus","index":1,"x":2,"y":10,"status":10}
//{"e":"changeUserStatus","index":1,"x":4,"y":14,"status":-1,"properties":[{"i":0,"x":1,"y":11,"type":0},{"i":13,"x":2,"y":11,"type":0},{"i":3,"x":3,"y":11,"type":1}]}
//{"e":"changeUserStatus","index":2,"x":4,"y":14,"status":-1,"properties":[{"i":42,"x":5,"y":11,"type":0},{"i":39,"x":6,"y":11,"type":1},{"i":29,"x":7,"y":11,"type":2}]}
type ChangeUserStatus struct {
	E          string        `json:"e"`
	Index      int           `json:"index"`
	X          int           `json:"x"`
	Y          int           `json:"y"`
	Status     int           `json:"status"`
	Properties []interface{} `json:"properties"`
}

//{"e":"gameOver","result":[{"_id":"63f08b2d8b2f2a4e2b30ddd9","index":0,"win":-1},{"_id":"63ec3b6a97611132f5ef439d","index":1,"win":0},{"_id":"63ec3b6a97611132f5ef43b0","index":2,"win":0}],"win":5,"lose":-2}

type GameOver struct {
	E      string `json:"e"`
	Result []GameOverResult `json:"result"`
	Win  int `json:"win"`
	Lose int `json:"lose"`
}

type GameOverResult struct {
	Id    string `json:"_id"`
	Index int    `json:"index"`
	Win   int    `json:"win"`
}