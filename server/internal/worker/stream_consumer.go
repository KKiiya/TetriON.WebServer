package worker

import (
	"context"
	"sync"
	"time"

	"TetriON.WebServer/server/internal/logging"
)

type StreamConsumer struct {
	name   string
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewStreamConsumer(name string) *StreamConsumer {
	if name == "" {
		name = "default"
	}
	return &StreamConsumer{name: name}
}

func (c *StreamConsumer) Start(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	c.cancel = cancel
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		logging.LogInfo("Stream consumer '%s' started", c.name)
		for {
			select {
			case <-ctx.Done():
				logging.LogInfo("Stream consumer '%s' stopped", c.name)
				return
			case <-ticker.C:
				// Placeholder polling tick to keep worker lifecycle in place.
			}
		}
	}()
}

func (c *StreamConsumer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
}
