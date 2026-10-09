package domain

import (
	"time"
)

// Metric encapsula la telemetría recolectada en cada inferencia o llamada de agente.
type Metric struct {
	ID               string    `json:"id"`
	InteractionID    string    `json:"interactionId,omitempty"`
	Provider         string    `json:"provider"`         // openai, gemini, anthropic, ollama, groq
	Model            string    `json:"model"`            // gpt-4o, gemini-1.5-flash, etc.
	TotalLatencyMs   int64     `json:"totalLatencyMs"`   // Latencia total en milisegundos
	TTFTMs           int64     `json:"ttftMs"`           // Time To First Token en milisegundos
	PromptTokens     int       `json:"promptTokens"`     // Tokens de entrada
	CompletionTokens int       `json:"completionTokens"` // Tokens generados
	TotalTokens      int       `json:"totalTokens"`      // Total de tokens
	FinishReason     string    `json:"finishReason"`     // stop, max_tokens, tool_calls
	Timestamp        time.Time `json:"timestamp"`
}

func NewMetric(id, provider, model string, totalLatencyMs, ttftMs int64, promptTokens, completionTokens int, finishReason string) *Metric {
	return &Metric{
		ID:               id,
		Provider:         provider,
		Model:            model,
		TotalLatencyMs:   totalLatencyMs,
		TTFTMs:           ttftMs,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
		FinishReason:     finishReason,
		Timestamp:        time.Now(),
	}
}
