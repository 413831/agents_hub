package domain

import (
	"errors"
	"time"
)

var (
	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidSessionTitle = errors.New("session title cannot be empty")
)

type SessionStatus string

const (
	SessionStatusActive   SessionStatus = "active"
	SessionStatusArchived SessionStatus = "archived"
)

// Session representa un hilo de conversación o contexto interactivo multi-agente.
type Session struct {
	ID        string        `json:"id"`
	ProjectID string        `json:"projectId"`
	Title     string        `json:"title"`
	AgentIDs  []string      `json:"agentIds"`
	Status    SessionStatus `json:"status"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func NewSession(id, projectID, title string, agentIDs []string) (*Session, error) {
	if title == "" {
		return nil, ErrInvalidSessionTitle
	}
	now := time.Now()
	if agentIDs == nil {
		agentIDs = make([]string, 0)
	}
	return &Session{
		ID:        id,
		ProjectID: projectID,
		Title:     title,
		AgentIDs:  agentIDs,
		Status:    SessionStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
