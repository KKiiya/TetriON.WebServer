package websocket

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan any
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
		case client := <-h.unregister:
			delete(h.clients, client.ID)
		case msg := <-h.broadcast:
			for _, c := range h.clients {
				c.Send <- msg
			}
		}
	}
}
