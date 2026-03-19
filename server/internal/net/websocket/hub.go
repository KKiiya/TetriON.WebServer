package websocket

import "sync"

type Hub struct {
	mu         sync.RWMutex
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan any
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client, 128),
		unregister: make(chan *Client, 128),
		broadcast:  make(chan any, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, exists := h.clients[client.ID]; exists {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			h.mu.RLock()
			for _, c := range h.clients {
				select {
				case c.Send <- msg:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(message any) {
	h.broadcast <- message
}

func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
