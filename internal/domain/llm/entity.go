package llm

import "time"

// ProviderType representa el identificador del motor LLM
type ProviderType string

const (
	ProviderOpenAI ProviderType = "openai"
	ProviderGemini ProviderType = "gemini"
	ProviderMock   ProviderType = "mock"
)

// Role indica el emisor del mensaje
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message representa un mensaje dentro del contexto conversacional
type Message struct {
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// CompletionRequest parámetros para solicitar inferencia
type CompletionRequest struct {
	Provider    ProviderType `json:"provider"`
	Model       string       `json:"model"`
	Messages    []Message    `json:"messages"`
	Temperature float32      `json:"temperature"`
	MaxTokens   int          `json:"max_tokens"`
}

// CompletionResponse resultado de la inferencia
type CompletionResponse struct {
	Provider     ProviderType `json:"provider"`
	Model        string       `json:"model"`
	Content      string       `json:"content"`
	PromptTokens int          `json:"prompt_tokens"`
	CompTokens   int          `json:"completion_tokens"`
	TotalTokens  int          `json:"total_tokens"`
	DurationMs   int64        `json:"duration_ms"`
}
