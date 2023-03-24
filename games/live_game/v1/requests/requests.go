package requests

type Common struct {
	R string `json:"r"`
}

//{"r":"login","_id":"63f08b2d8b2f2a4e2b30ddd9","id":100001,"token":"12312313"}
type Login struct {
	R     string `json:"r"`
	Id    string `json:"_id"`
	Id1   int    `json:"id"`
	Token string `json:"token"`
}

//{"r":"matching","role":300,"bombColor":0}
type Matching struct {
	R         string `json:"r"`
	Role      int    `json:"role"`
	BombColor int    `json:"bombColor"`
}

//{"r":"moveStop"}
type MoveStop struct {
	R         string `json:"r"`
}

//{ r: 'setPos', x: x, y: y, n: this._lastMove };

type SetPos struct {
	R         string `json:"r"`
	X int `json:"x"`
	Y int `json:"y"`
	N int `json:"n"`
}

func GetRoomMessage() interface{}{
	return Common{
		R: "getRoomMessage",
	}
}