package broker

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	domainTask "proyecto_ia/internal/domain/task"
)

type InMemoryJobRepository struct {
	mu   sync.RWMutex
	jobs map[string]domainTask.Job
}

func NewInMemoryJobRepository() *InMemoryJobRepository {
	return &InMemoryJobRepository{
		jobs: make(map[string]domainTask.Job),
	}
}

func (r *InMemoryJobRepository) Save(ctx context.Context, job domainTask.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryJobRepository) GetByID(ctx context.Context, id string) (*domainTask.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", id)
	}
	return &job, nil
}

func (r *InMemoryJobRepository) ListRecent(ctx context.Context, limit int) ([]domainTask.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domainTask.Job, 0, len(r.jobs))
	for _, j := range r.jobs {
		list = append(list, j)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.After(list[j].CreatedAt)
	})

	if len(list) > limit {
		list = list[:limit]
	}
	return list, nil
}

func (r *InMemoryJobRepository) UpdateStatus(ctx context.Context, id string, status domainTask.TaskStatus, progress int, result []byte, errStr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[id]
	if !exists {
		return fmt.Errorf("job not found: %s", id)
	}

	job.Status = status
	job.Progress = progress
	if len(result) > 0 {
		job.Result = result
	}
	job.Error = errStr
	job.UpdatedAt = time.Now()

	r.jobs[id] = job
	return nil
}
