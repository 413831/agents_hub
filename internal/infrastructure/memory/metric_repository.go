package memory

import (
	"context"
	"sort"
	"sync"

	"proyecto_ia/internal/domain"
)

type InMemoryMetricRepository struct {
	mu      sync.RWMutex
	metrics []domain.Metric
}

func NewInMemoryMetricRepository() *InMemoryMetricRepository {
	return &InMemoryMetricRepository{
		metrics: make([]domain.Metric, 0),
	}
}

func (r *InMemoryMetricRepository) Save(ctx context.Context, m *domain.Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics = append(r.metrics, *m)
	return nil
}

func (r *InMemoryMetricRepository) ListRecent(ctx context.Context, limit int) ([]domain.Metric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.metrics) == 0 {
		return []domain.Metric{}, nil
	}
	out := make([]domain.Metric, len(r.metrics))
	copy(out, r.metrics)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Timestamp.After(out[j].Timestamp)
	})
	if limit > 0 && len(out) > limit {
		return out[:limit], nil
	}
	return out, nil
}
