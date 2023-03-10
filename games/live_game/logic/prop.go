package logic

type Prop struct {
	Id int
	X int
	Y int
	Type int
	RoomId int
}
type PropMsg struct {
	Id int `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Type int `json:"type"`
	Room int `json:"room"`
}

func (p Prop) GetMsg() PropMsg {
	return PropMsg{
		Id:   p.Id,
		X:    p.X,
		Y:    p.Y,
		Type: p.Type,
		Room: p.RoomId,
	}
}