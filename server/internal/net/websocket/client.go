package websocket

import "github.com/coder/websocket"

type Client struct {
	Conn *websocket.Conn
	ID   string
	Send chan any
}
