package gameserver

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var ErrNoAvailableServer = errors.New("no available game server")

type Node struct {
	ID          string    `json:"id"`
	Address     string    `json:"address"`
	Capacity    int       `json:"capacity"`
	CurrentLoad int       `json:"current_load"`
	LastSeen    time.Time `json:"last_seen"`
}

type Registry struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

func NewRegistry() *Registry {
	return &Registry{nodes: make(map[string]*Node)}
}

func (r *Registry) Upsert(node Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	node.LastSeen = time.Now()
	n := node
	r.nodes[node.ID] = &n
}

func (r *Registry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.nodes, id)
}

func (r *Registry) Heartbeat(id string, currentLoad int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if node, ok := r.nodes[id]; ok {
		node.CurrentLoad = currentLoad
		node.LastSeen = time.Now()
	}
}

func (r *Registry) List() []Node {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Node, 0, len(r.nodes))
	for _, n := range r.nodes {
		out = append(out, *n)
	}
	return out
}

func (r *Registry) SelectLeastLoaded() (Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	candidates := make([]*Node, 0, len(r.nodes))
	for _, n := range r.nodes {
		if n.Capacity > 0 && n.CurrentLoad < n.Capacity {
			candidates = append(candidates, n)
		}
	}

	if len(candidates) == 0 {
		return Node{}, ErrNoAvailableServer
	}

	sort.Slice(candidates, func(i, j int) bool {
		li := float64(candidates[i].CurrentLoad) / float64(candidates[i].Capacity)
		lj := float64(candidates[j].CurrentLoad) / float64(candidates[j].Capacity)
		if li == lj {
			return candidates[i].LastSeen.After(candidates[j].LastSeen)
		}
		return li < lj
	})

	return *candidates[0], nil
}

func (r *Registry) PruneStale(maxAge time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for id, n := range r.nodes {
		if now.Sub(n.LastSeen) > maxAge {
			delete(r.nodes, id)
		}
	}
}
