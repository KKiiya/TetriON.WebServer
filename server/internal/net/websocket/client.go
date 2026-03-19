package websocket

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Client struct {
	Conn *websocket.Conn
	ID   string
	Send chan any
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		ID:   id,
		Send: make(chan any, 64),
	}
}

func (c *Client) ReadPump(ctx context.Context, onMessage func(any)) {
	for {
		var payload any
		if err := wsjson.Read(ctx, c.Conn, &payload); err != nil {
			return
		}
		onMessage(payload)
	}
}

func (c *Client) WritePump(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.Send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			_ = wsjson.Write(writeCtx, c.Conn, msg)
			cancel()
		}
	}
}
