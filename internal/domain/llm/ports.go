package llm

import "context"

// ChunkCallback función callback para streaming de tokens
type ChunkCallback func(chunk string) error

// LLMProvider puerto primario de dominio para interactuar con proveedores LLM
type LLMProvider interface {
	Type() ProviderType
	SupportedModels() []string
	GenerateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	StreamCompletion(ctx context.Context, req CompletionRequest, onChunk ChunkCallback) error
}

// LLMRepository puerto secundario opcional para guardar historial de prompts
type LLMHistoryRepository interface {
	SaveMessage(ctx context.Context, sessionID string, msg Message) error
	GetHistory(ctx context.Context, sessionID string) ([]Message, error)
}
