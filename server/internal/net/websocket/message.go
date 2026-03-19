package websocket

type Message struct {
	Type      string         `json:"type"`
	Source    string         `json:"source,omitempty"`
	Timestamp int64          `json:"timestamp"`
	Payload   map[string]any `json:"payload,omitempty"`
}
