package game

import "encoding/json"

type Message struct {
	Event string `json:"e"`
	Info string `json:"info,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

type Handler interface {
	HandleMessage(pipe *Pipe, msg *Message)
}