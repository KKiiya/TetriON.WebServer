package matchmaking

import (
	"context"

	redisnet "TetriON.WebServer/server/internal/net/redis"
)

type Manager struct {
	queueName string
}

func NewManager(queueName string) *Manager {
	if queueName == "" {
		queueName = "default"
	}
	return &Manager{queueName: queueName}
}

func (m *Manager) Enqueue(ctx context.Context, userID string, skill int) error {
	return redisnet.EnqueuePlayer(ctx, m.queueName, userID, skill)
}

func (m *Manager) Dequeue(ctx context.Context, userID string) error {
	return redisnet.RemovePlayerFromQueue(ctx, m.queueName, userID)
}

func (m *Manager) Candidates(ctx context.Context, limit int64) ([]string, error) {
	return redisnet.PeekPlayers(ctx, m.queueName, limit)
}

func (m *Manager) Size(ctx context.Context) (int64, error) {
	return redisnet.QueueSize(ctx, m.queueName)
}
