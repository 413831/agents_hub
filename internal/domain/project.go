package domain

import (
	"errors"
	"time"
)

var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrInvalidProjectName  = errors.New("project name cannot be empty")
)

// Project representa un espacio de trabajo o contexto para agrupar sesiones y agentes.
type Project struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	DefaultModel string           `json:"defaultModel,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

func NewProject(id, name, description string) (*Project, error) {
	if name == "" {
		return nil, ErrInvalidProjectName
	}
	now := time.Now()
	return &Project{
		ID:          id,
		Name:        name,
		Description: description,
		Metadata:    make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
