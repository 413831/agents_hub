package domain

import (
	"context"
)

// ProjectRepository define el contrato para persistencia y consulta de Proyectos.
type ProjectRepository interface {
	Save(ctx context.Context, project *Project) error
	FindByID(ctx context.Context, id string) (*Project, error)
	List(ctx context.Context) ([]Project, error)
	Delete(ctx context.Context, id string) error
}

// SessionRepository define el contrato para persistencia y consulta de Sesiones.
type SessionRepository interface {
	Save(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id string) (*Session, error)
	ListByProject(ctx context.Context, projectID string) ([]Session, error)
	Delete(ctx context.Context, id string) error
}

// MessageRepository define el contrato para persistencia y consulta de Mensajes.
type MessageRepository interface {
	Save(ctx context.Context, message *Message) error
	ListBySession(ctx context.Context, sessionID string) ([]Message, error)
	DeleteBySession(ctx context.Context, sessionID string) error
}

// AgentRepository define el contrato para persistencia y consulta de Agentes.
type AgentRepository interface {
	Save(ctx context.Context, agent *Agent) error
	FindByID(ctx context.Context, id string) (*Agent, error)
	List(ctx context.Context) ([]Agent, error)
	Delete(ctx context.Context, id string) error
}

// ConnectorRepository define el contrato para persistencia y consulta de Conectores (MCP, Skills, Plugins).
type ConnectorRepository interface {
	Save(ctx context.Context, connector *Connector) error
	FindByID(ctx context.Context, id string) (*Connector, error)
	List(ctx context.Context) ([]Connector, error)
	UpdateStatus(ctx context.Context, id string, status ConnectorStatus) error
	Delete(ctx context.Context, id string) error
}

// MetricRepository define el contrato para almacenamiento y agregación de telemetría.
type MetricRepository interface {
	Save(ctx context.Context, metric *Metric) error
	ListRecent(ctx context.Context, limit int) ([]Metric, error)
}
