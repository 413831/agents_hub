package domain

import (
	"errors"
	"time"
)

var (
	ErrConnectorNotFound = errors.New("connector not found")
)

type ConnectorType string

const (
	ConnectorTypeMCP    ConnectorType = "mcp"
	ConnectorTypeSkill  ConnectorType = "skill"
	ConnectorTypePlugin ConnectorType = "plugin"
)

type ConnectorStatus string

const (
	ConnectorStatusConnected    ConnectorStatus = "connected"
	ConnectorStatusDisconnected ConnectorStatus = "disconnected"
	ConnectorStatusDisabled     ConnectorStatus = "disabled"
)

// Connector representa una conexión a una tool, servidor MCP o plugin externo.
type Connector struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        ConnectorType     `json:"type"`
	Description string            `json:"description"`
	Status      ConnectorStatus   `json:"status"`
	Endpoint    string            `json:"endpoint,omitempty"` // URL o comando ejecutable MCP
	Tools       []string          `json:"tools"`              // Herramientas expuestas (ej: "read_file", "git_status")
	Config      map[string]string `json:"config,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

func NewConnector(id, name string, connType ConnectorType, description, endpoint string, tools []string) *Connector {
	now := time.Now()
	if tools == nil {
		tools = make([]string, 0)
	}
	return &Connector{
		ID:          id,
		Name:        name,
		Type:        connType,
		Description: description,
		Status:      ConnectorStatusConnected,
		Endpoint:    endpoint,
		Tools:       tools,
		Config:      make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
