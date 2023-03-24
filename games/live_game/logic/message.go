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

type LoginData struct {
	Id string `json:"id"`
	Token string `json:"token"`
}

type LoadProcessData struct {
	Id int `json:"id"`
	Process int `json:"process"`
}

type JoinData struct {
	Role      int    `json:"role"`
	BombColor int    `json:"bombColor"`
}

type JoinSuccessData struct {
	E    string `json:"e"`
	Self *PlayerData `json:"self"`
	Users []*PlayerData `json:"users"`
}

type RoundStartData struct {
	RoomId int `json:"roomId"`
	Type int `json:"type"`
	RoomType int `json:"roomType"`
	Status int `json:"status"`
	Player []PlayerData `json:"users"`
	Map maps.MapData `json:"map"`
	Bubbles []int `json:"bombs"`
	Props []PropData `json:"props"`
}

type PropDisappearData struct {
	Id int `json:"id"`
	PickUp bool `json:"pick_up"`
	Player string `json:"player"`
}