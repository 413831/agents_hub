package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"proyecto_ia/internal/domain"
)

type InMemoryMessageRepository struct {
	mu       sync.RWMutex
	messages map[string][]domain.Message // key: sessionId
}

func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	repo := &InMemoryMessageRepository{
		messages: make(map[string][]domain.Message),
	}
	repo.seed()
	return repo
}

func (r *InMemoryMessageRepository) seed() {
	now := time.Now()
	// Mensajes para sess-1
	m1 := domain.Message{
		ID:        "m1-user",
		SessionID: "sess-1",
		Role:      domain.RoleUser,
		Sender:    "Developer",
		Content:   "Hola, ¿cómo estructuramos las capas de Clean Architecture para el Hub de Agentes?",
		CreatedAt: now.Add(-30 * time.Minute),
	}
	m2 := domain.Message{
		ID:        "m2-agent",
		SessionID: "sess-1",
		AgentID:   "agent-arch-1",
		Role:      domain.RoleAssistant,
		Sender:    "Architect Core",
		Content:   "La estructura recomendada sigue la regla de dependencias concéntricas:\n\n1. **domain**: Entidades puras y puertos (interfaces).\n2. **usecase**: Casos de uso de negocio sin acoplarse a frameworks.\n3. **infrastructure**: Implementaciones concretas de proveedores LLM, brokers y repositorios.\n4. **interfaces**: Controladores y adaptadores de entrada (Wails bindings, CLI).",
		Telemetry: &domain.Metric{
			ID:               "met-seed-1",
			Provider:         "google",
			Model:            "gemini-1.5-flash",
			TotalLatencyMs:   342,
			TTFTMs:           110,
			PromptTokens:     42,
			CompletionTokens: 118,
			TotalTokens:      160,
			FinishReason:     "stop",
			Timestamp:        now.Add(-30 * time.Minute),
		},
		CreatedAt: now.Add(-29 * time.Minute),
	}
	r.messages["sess-1"] = []domain.Message{m1, m2}
}

func (r *InMemoryMessageRepository) Save(ctx context.Context, m *domain.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages[m.SessionID] = append(r.messages[m.SessionID], *m)
	return nil
}

func (r *InMemoryMessageRepository) ListBySession(ctx context.Context, sessionID string) ([]domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, exists := r.messages[sessionID]
	if !exists {
		return []domain.Message{}, nil
	}
	out := make([]domain.Message, len(list))
	copy(out, list)
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (r *InMemoryMessageRepository) DeleteBySession(ctx context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.messages, sessionID)
	return nil
}
