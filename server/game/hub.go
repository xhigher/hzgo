package game

import (
	"encoding/json"
	"time"
)

type Hub struct {
	// Registered clients.
	pipes map[*Pipe]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Pipe

	// Unregister requests from clients.
	unregister chan *Pipe

	handler Handler

	heartbeatTimer *time.Ticker
}

func newHub(handler Handler) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Pipe),
		unregister: make(chan *Pipe),
		pipes:    make(map[*Pipe]bool),
		handler: handler,
		heartbeatTimer: time.NewTicker(1*time.Second),
	}
}

func (h *Hub) run() {
	hbMsg,_ := json.Marshal(&Message{
		Event: "heartbeat",
	})
	for {
		select {
		case <-h.heartbeatTimer.C:
			h.broadcastMsg(hbMsg)
		case p := <-h.register:
			h.addPipe(p)
		case p := <-h.unregister:
			h.removePipe(p)
		case msg := <-h.broadcast:
			h.broadcastMsg(msg)
		}
	}
}

func (h *Hub) addPipe(p *Pipe){
	h.pipes[p] = true
}

func (h *Hub) removePipe(p *Pipe){
	if _, ok := h.pipes[p]; ok {
		delete(h.pipes, p)
	}
}

func (h *Hub) broadcastMsg(msg []byte){
	for p := range h.pipes {
		p.writeMessage(msg)
	}
}


