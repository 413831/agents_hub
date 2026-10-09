package domain

import (
	"errors"
	"time"
)

var (
	ErrAgentNotFound    = errors.New("agent not found")
	ErrInvalidAgentName = errors.New("agent name cannot be empty")
)

// Agent representa un agente autónomo o asistente configurado con rol, system prompt y modelo asignado.
type Agent struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	Description  string    `json:"description"`
	SystemPrompt string    `json:"systemPrompt"`
	Provider     string    `json:"provider"`    // openai, gemini, anthropic, ollama, groq
	Model        string    `json:"model"`       // gpt-4o, gemini-1.5-flash, claude-3-5-sonnet, llama3
	Temperature  float64   `json:"temperature"` // 0.0 - 1.0
	MaxTokens    int       `json:"maxTokens"`
	Tools        []string  `json:"tools"`       // IDs de conectores o tools permitidas
	Avatar       string    `json:"avatar,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewAgent(id, name, role, description, systemPrompt, provider, model string) (*Agent, error) {
	if name == "" {
		return nil, ErrInvalidAgentName
	}
	now := time.Now()
	return &Agent{
		ID:           id,
		Name:         name,
		Role:         role,
		Description:  description,
		SystemPrompt: systemPrompt,
		Provider:     provider,
		Model:        model,
		Temperature:  0.7,
		MaxTokens:    2048,
		Tools:        make([]string, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
