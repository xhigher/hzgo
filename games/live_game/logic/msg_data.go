package logic

import (
	"encoding/json"
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

type LoginData struct {
	Id string `json:"id"`
	Token string `json:"token"`
}

type LoadProcessData struct {
	Id int `json:"id"`
	Process int `json:"process"`
}
