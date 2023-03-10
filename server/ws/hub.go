package ws

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
}

func newHub(handler Handler) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Pipe),
		unregister: make(chan *Pipe),
		pipes:    make(map[*Pipe]bool),
		handler: handler,
	}
}

func (h *Hub) run() {
	for {
		select {
		case pipe := <-h.register:
			h.pipes[pipe] = true
		case pipe := <-h.unregister:
			if _, ok := h.pipes[pipe]; ok {
				delete(h.pipes, pipe)
				close(pipe.send)
			}
		case msg := <-h.broadcast:
			for pipe := range h.pipes {
				select {
				case pipe.send <- msg:
				default:
					close(pipe.send)
					delete(h.pipes, pipe)
				}
			}
		}
	}
}


