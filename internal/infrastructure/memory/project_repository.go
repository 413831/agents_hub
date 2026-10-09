package memory

import (
	"context"
	"sync"
	"time"

	"proyecto_ia/internal/domain"
)

type InMemoryProjectRepository struct {
	mu       sync.RWMutex
	projects map[string]domain.Project
}

func NewInMemoryProjectRepository() *InMemoryProjectRepository {
	repo := &InMemoryProjectRepository{
		projects: make(map[string]domain.Project),
	}
	repo.seed()
	return repo
}

func (r *InMemoryProjectRepository) seed() {
	now := time.Now()
	p1 := domain.Project{
		ID:           "proj-default-1",
		Name:         "AI Agent Laboratory",
		Description:  "Experimentación multi-proveedor con OpenAI, Gemini y Llama local.",
		DefaultModel: "gemini-1.5-flash",
		Metadata:     map[string]string{"env": "development", "version": "1.0"},
		CreatedAt:    now.Add(-2 * time.Hour),
		UpdatedAt:    now,
	}
	p2 := domain.Project{
		ID:           "proj-default-2",
		Name:         "Code Refactor & MCP Tools",
		Description:  "Pipeline de desarrollo y análisis estático con conectores MCP locales.",
		DefaultModel: "gpt-4o",
		Metadata:     map[string]string{"env": "local", "tools": "git,files"},
		CreatedAt:    now.Add(-1 * time.Hour),
		UpdatedAt:    now,
	}
	r.projects[p1.ID] = p1
	r.projects[p2.ID] = p2
}

func (r *InMemoryProjectRepository) Save(ctx context.Context, p *domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projects[p.ID] = *p
	return nil
}

func (r *InMemoryProjectRepository) FindByID(ctx context.Context, id string) (*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.projects[id]
	if !exists {
		return nil, domain.ErrProjectNotFound
	}
	return &p, nil
}

func (r *InMemoryProjectRepository) List(ctx context.Context) ([]domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.Project, 0, len(r.projects))
	for _, p := range r.projects {
		list = append(list, p)
	}
	return list, nil
}

func (r *InMemoryProjectRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.projects, id)
	return nil
}
