package memory

import (
	"context"
	"sync"
	"time"

	"proyecto_ia/internal/domain"
)

type InMemoryAgentRepository struct {
	mu     sync.RWMutex
	agents map[string]domain.Agent
}

func NewInMemoryAgentRepository() *InMemoryAgentRepository {
	repo := &InMemoryAgentRepository{
		agents: make(map[string]domain.Agent),
	}
	repo.seed()
	return repo
}

func (r *InMemoryAgentRepository) seed() {
	now := time.Now()
	a1 := domain.Agent{
		ID:           "agent-arch-1",
		Name:         "Architect Core",
		Role:         "Principal Software Architect",
		Description:  "Especialista en diseño de sistemas distribuidos, Clean Architecture y DDD.",
		SystemPrompt: "Eres un arquitecto de software senior enfocado en patrones limpios y rendimiento en Golang y TypeScript.",
		Provider:     "google",
		Model:        "gemini-1.5-flash",
		Temperature:  0.4,
		MaxTokens:    4096,
		Tools:        []string{"conn-git-1", "conn-mcp-fs"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	a2 := domain.Agent{
		ID:           "agent-code-1",
		Name:         "Code Reviewer",
		Role:         "Senior Staff Engineer",
		Description:  "Analiza código fuente, detección de vulnerabilidades y optimización de concurrencia.",
		SystemPrompt: "Eres un revisor de código estricto especializado en buenas prácticas, tests unitarios y seguridad.",
		Provider:     "openai",
		Model:        "gpt-4o",
		Temperature:  0.2,
		MaxTokens:    4096,
		Tools:        []string{"conn-git-1"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	a3 := domain.Agent{
		ID:           "agent-local-1",
		Name:         "Local Llama Copilot",
		Role:         "Edge Assistant",
		Description:  "Ejecución en hardware local offline vía Ollama/Groq.",
		SystemPrompt: "Eres un asistente de terminal conciso y rápido que opera localmente.",
		Provider:     "ollama",
		Model:        "llama3:8b",
		Temperature:  0.7,
		MaxTokens:    2048,
		Tools:        []string{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	r.agents[a1.ID] = a1
	r.agents[a2.ID] = a2
	r.agents[a3.ID] = a3
}

func (r *InMemoryAgentRepository) Save(ctx context.Context, a *domain.Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[a.ID] = *a
	return nil
}

func (r *InMemoryAgentRepository) FindByID(ctx context.Context, id string) (*domain.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, exists := r.agents[id]
	if !exists {
		return nil, domain.ErrAgentNotFound
	}
	return &a, nil
}

func (r *InMemoryAgentRepository) List(ctx context.Context) ([]domain.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.Agent, 0, len(r.agents))
	for _, a := range r.agents {
		list = append(list, a)
	}
	return list, nil
}

func (r *InMemoryAgentRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.agents, id)
	return nil
}
