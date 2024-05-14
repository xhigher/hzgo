package logic

import (
	"encoding/json"
	"github.com/xhigher/hzgo/games/live_game/maps"
	"github.com/xhigher/hzgo/logger"
)

func encodeMsgData(data interface{}) json.RawMessage {
	bs, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("error: %v", err)
	}
	return bs
}

func decodeMsgData(raw json.RawMessage, data interface{}) error {
	err := json.Unmarshal(raw, data)
	if err != nil {
		logger.Errorf("error: %v", err)
	}
	return err
}

type PlayerId struct {
	Id string `json:"id"`
}

type PlayerIndex struct {
	I int `json:"i"`
}

type LoginData struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

type LoadProcessData struct {
	Id      int `json:"id"`
	Process int `json:"process"`
}

type JoinData struct {
	Skin        int `json:"role"`
	BubbleColor int `json:"bombColor"`
}

type JoinSuccessData struct {
	E     string        `json:"e"`
	Self  *PlayerData   `json:"self"`
	Users []*PlayerData `json:"users"`
}

type RoundStartData struct {
	RoomId   int          `json:"roomId"`
	Type     int          `json:"type"`
	RoomType int          `json:"roomType"`
	Status   int          `json:"status"`
	Player   []PlayerData `json:"users"`
	Map      maps.MapData `json:"map"`
	Bubbles  []int        `json:"bombs"`
	Props    []PropData   `json:"props"`
}

type PropDisappearData struct {
	Id     int    `json:"i"`
	PickUp bool   `json:"isEat"`
	Player string `json:"userIndex"`
}

type MoveData struct {
	I int `json:"i"`
	X int `json:"x"`
	Y int `json:"y"`
	T int `json:"t"`
	N int `json:"n"`
}

type RoundOverData struct {
	Result []RoundOverResult `json:"result"`
	Win    int               `json:"win"`
	Lose   int               `json:"lose"`
}

type RoundOverResult struct {
	Player string `json:"id"`
	Index  int    `json:"index"`
	Win    int    `json:"win"`
}

type ChangeUserStatusData struct {
	Id     string `json:"id"`
	Index  int    `json:"index"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Status int    `json:"status"`
}
