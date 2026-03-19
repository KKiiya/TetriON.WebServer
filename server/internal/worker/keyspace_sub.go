package worker

import (
	"context"
	"sync"

	"TetriON.WebServer/server/internal/logging"
	redisnet "TetriON.WebServer/server/internal/net/redis"
)

type KeyspaceSubscriber struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewKeyspaceSubscriber() *KeyspaceSubscriber {
	return &KeyspaceSubscriber{}
}

func (s *KeyspaceSubscriber) Start(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	s.cancel = cancel
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		logging.LogInfo("Starting Redis pub/sub subscriber worker")
		redisnet.SubscribeMessages(ctx)
		logging.LogInfo("Redis pub/sub subscriber worker stopped")
	}()
}

func (s *KeyspaceSubscriber) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
