package ws

import "encoding/json"

type Message struct {
	Event string `json:"event"`
	Info string `json:"info"`
	Data json.RawMessage `json:"data"`
}

type Handler interface {
	HandleMessage(pipe *Pipe, msg *Message)
}