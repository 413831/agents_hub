package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidMessageContent = errors.New("message content cannot be empty")
)

type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
	RoleTool      MessageRole = "tool"
)

// Message representa un mensaje dentro de una sesión de conversación.
type Message struct {
	ID        string       `json:"id"`
	SessionID string       `json:"sessionId"`
	AgentID   string       `json:"agentId,omitempty"` // Si es emitido por un agente
	Role      MessageRole  `json:"role"`
	Sender    string       `json:"sender"`            // Nombre legible (ej: "Usuario", "Architect", "Gemini")
	Content   string       `json:"content"`
	Telemetry *Metric      `json:"telemetry,omitempty"`
	CreatedAt time.Time    `json:"createdAt"`
}

func NewMessage(id, sessionID, agentID string, role MessageRole, sender, content string, telemetry *Metric) (*Message, error) {
	if content == "" {
		return nil, ErrInvalidMessageContent
	}
	return &Message{
		ID:        id,
		SessionID: sessionID,
		AgentID:   agentID,
		Role:      role,
		Sender:    sender,
		Content:   content,
		Telemetry: telemetry,
		CreatedAt: time.Now(),
	}, nil
}
