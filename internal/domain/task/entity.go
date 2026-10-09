package task

import (
	"encoding/json"
	"time"
)

// TaskStatus representa los estados del ciclo de vida de un Job
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

// TaskType identifica el tipo de trabajo asíncrono
type TaskType string

const (
	TypeLLMInference TaskType = "task:llm:inference"
	TypeGitClone     TaskType = "task:git:clone"
	TypeGitAnalysis  TaskType = "task:git:analysis"
)

// Job entidad que encapsula una tarea en segundo plano
type Job struct {
	ID        string          `json:"id"`
	Type      TaskType        `json:"type"`
	Status    TaskStatus      `json:"status"`
	Progress  int             `json:"progress"` // 0 - 100%
	Payload   json.RawMessage `json:"payload"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     string          `json:"error,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
