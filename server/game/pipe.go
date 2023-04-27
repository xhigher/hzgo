package game

import (
	"bytes"
	"encoding/json"
	"github.com/xhigher/hzgo/logger"
	"log"
	"sync"
	"time"

	"github.com/hertz-contrib/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Pipe struct {
	hub       *Hub
	tag       string
	conn      *websocket.Conn
	send      chan []byte
	closeLock sync.Once
}

func NewPipe(hub *Hub, conn *websocket.Conn) *Pipe {
	p := &Pipe{
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, 256),
		closeLock: sync.Once{},
	}
	hub.register <- p
	return p
}

func (p *Pipe) readPump() {
	defer func() {
		p.Close()
		p.hub.unregister <- p
	}()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, mbs, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		mbs = bytes.TrimSpace(bytes.Replace(mbs, newline, space, -1))
		p.handleMessage(mbs)
	}
}

func (p *Pipe) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.Close()
	}()
	for {
		select {
		case mbs, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(mbs)

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (p *Pipe) SetTag(tag string) {
	p.tag = tag
}

func (p *Pipe) GetTag() string {
	return p.tag
}

func (p *Pipe) SendMessage(msg *Message) {
	mbs, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	logger.Infof("SendMessage: %v", string(mbs))
	p.writeMessage(mbs)
}

func (p *Pipe) writeMessage(mbs []byte) {
	defer func() {
		recover()
	}()
	if p.send != nil {
		select {
		case p.send <- mbs:
		default:
			logger.Errorf("connection pipe blocked: %v", p.conn.RemoteAddr())
		}

	}
}

func (p *Pipe) handleMessage(mbs []byte) {
	var msg *Message
	err := json.Unmarshal(mbs, &msg)
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	if p.hub.handler != nil {
		p.hub.handler.HandleMessage(p, msg)
	}
}

func (p *Pipe) Close() {
	p.closeLock.Do(func() {
		close(p.send)
		p.conn.Close()
	})
}
