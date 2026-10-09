package memory

import (
	"context"
	"sync"
	"time"

	"proyecto_ia/internal/domain"
)

type InMemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]domain.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	repo := &InMemorySessionRepository{
		sessions: make(map[string]domain.Session),
	}
	repo.seed()
	return repo
}

func (r *InMemorySessionRepository) seed() {
	now := time.Now()
	s1 := domain.Session{
		ID:        "sess-1",
		ProjectID: "proj-default-1",
		Title:     "Arquitectura DDD y Clean Architecture",
		AgentIDs:  []string{"agent-arch-1"},
		Status:    domain.SessionStatusActive,
		CreatedAt: now.Add(-40 * time.Minute),
		UpdatedAt: now.Add(-5 * time.Minute),
	}
	s2 := domain.Session{
		ID:        "sess-2",
		ProjectID: "proj-default-1",
		Title:     "Benchmark Latencia Gemini Flash vs OpenAI",
		AgentIDs:  []string{"agent-arch-1", "agent-code-1"},
		Status:    domain.SessionStatusActive,
		CreatedAt: now.Add(-20 * time.Minute),
		UpdatedAt: now.Add(-2 * time.Minute),
	}
	s3 := domain.Session{
		ID:        "sess-3",
		ProjectID: "proj-default-2",
		Title:     "Auditoría de Seguridad y MCP Tools",
		AgentIDs:  []string{"agent-code-1"},
		Status:    domain.SessionStatusActive,
		CreatedAt: now.Add(-10 * time.Minute),
		UpdatedAt: now,
	}
	r.sessions[s1.ID] = s1
	r.sessions[s2.ID] = s2
	r.sessions[s3.ID] = s3
}

func (r *InMemorySessionRepository) Save(ctx context.Context, s *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s.ID] = *s
	return nil
}

func (r *InMemorySessionRepository) FindByID(ctx context.Context, id string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, exists := r.sessions[id]
	if !exists {
		return nil, domain.ErrSessionNotFound
	}
	return &s, nil
}

func (r *InMemorySessionRepository) ListByProject(ctx context.Context, projectID string) ([]domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.Session, 0)
	for _, s := range r.sessions {
		if projectID == "" || s.ProjectID == projectID {
			list = append(list, s)
		}
	}
	return list, nil
}

func (r *InMemorySessionRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id)
	return nil
}
