package memory

import (
	"context"
	"sync"
	"time"

	"proyecto_ia/internal/domain"
)

type InMemoryConnectorRepository struct {
	mu         sync.RWMutex
	connectors map[string]domain.Connector
}

func NewInMemoryConnectorRepository() *InMemoryConnectorRepository {
	repo := &InMemoryConnectorRepository{
		connectors: make(map[string]domain.Connector),
	}
	repo.seed()
	return repo
}

func (r *InMemoryConnectorRepository) seed() {
	now := time.Now()
	c1 := domain.Connector{
		ID:          "conn-git-1",
		Name:        "Git Version Control MCP",
		Type:        domain.ConnectorTypeMCP,
		Description: "Permite inspeccionar diffs, commits, branches y staging nativo de Git.",
		Status:      domain.ConnectorStatusConnected,
		Endpoint:    "mcp://local/git-service",
		Tools:       []string{"git_status", "git_diff", "git_commit", "git_log"},
		Config:      map[string]string{"working_dir": "./"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	c2 := domain.Connector{
		ID:          "conn-mcp-fs",
		Name:        "Filesystem & AST MCP Server",
		Type:        domain.ConnectorTypeMCP,
		Description: "Exploración de archivos, parseo AST y búsqueda semántica de símbolos.",
		Status:      domain.ConnectorStatusConnected,
		Endpoint:    "mcp://local/filesystem",
		Tools:       []string{"read_file", "write_file", "search_symbols", "list_directory"},
		Config:      map[string]string{"read_only": "false"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	c3 := domain.Connector{
		ID:          "conn-skill-docker",
		Name:        "Docker & Container Skill",
		Type:        domain.ConnectorTypeSkill,
		Description: "Manejo de contenedores, docker-compose y logs de servicios.",
		Status:      domain.ConnectorStatusDisabled,
		Endpoint:    "skill://system/docker",
		Tools:       []string{"docker_ps", "docker_logs", "docker_restart"},
		Config:      map[string]string{"socket": "/var/run/docker.sock"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	c4 := domain.Connector{
		ID:          "conn-plugin-db",
		Name:        "PostgreSQL Inspector Plugin",
		Type:        domain.ConnectorTypePlugin,
		Description: "Introspección de esquemas relacionales, índices y EXPLAIN ANALYZE.",
		Status:      domain.ConnectorStatusDisconnected,
		Endpoint:    "plugin://postgres/introspect",
		Tools:       []string{"db_schema", "db_query_plan"},
		Config:      map[string]string{"host": "localhost:5432"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.connectors[c1.ID] = c1
	r.connectors[c2.ID] = c2
	r.connectors[c3.ID] = c3
	r.connectors[c4.ID] = c4
}

func (r *InMemoryConnectorRepository) Save(ctx context.Context, c *domain.Connector) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connectors[c.ID] = *c
	return nil
}

func (r *InMemoryConnectorRepository) FindByID(ctx context.Context, id string) (*domain.Connector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.connectors[id]
	if !exists {
		return nil, domain.ErrConnectorNotFound
	}
	return &c, nil
}

func (r *InMemoryConnectorRepository) List(ctx context.Context) ([]domain.Connector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.Connector, 0, len(r.connectors))
	for _, c := range r.connectors {
		list = append(list, c)
	}
	return list, nil
}

func (r *InMemoryConnectorRepository) UpdateStatus(ctx context.Context, id string, status domain.ConnectorStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, exists := r.connectors[id]
	if !exists {
		return domain.ErrConnectorNotFound
	}
	c.Status = status
	c.UpdatedAt = time.Now()
	r.connectors[id] = c
	return nil
}

func (r *InMemoryConnectorRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.connectors, id)
	return nil
}
