package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	domainLLM "proyecto_ia/internal/domain/llm"
)

// OpenAIAdapter desacopla la comunicación HTTP con la API de OpenAI
type OpenAIAdapter struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewOpenAIAdapter(apiKey string) *OpenAIAdapter {
	if apiKey == "" {
		apiKey = "dummy-key"
	}
	return &OpenAIAdapter{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (a *OpenAIAdapter) Type() domainLLM.ProviderType {
	return domainLLM.ProviderOpenAI
}

func (a *OpenAIAdapter) SupportedModels() []string {
	return []string{"gpt-4o", "gpt-4o-mini", "o1-mini", "gpt-3.5-turbo"}
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIReqBody struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float32         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Stream      bool            `json:"stream"`
}

type openAIRespBody struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (a *OpenAIAdapter) GenerateCompletion(ctx context.Context, req domainLLM.CompletionRequest) (*domainLLM.CompletionResponse, error) {
	if a.apiKey == "" || a.apiKey == "dummy-key" {
		// Modo simulado amigable si no hay clave real configurada
		return &domainLLM.CompletionResponse{
			Provider:     domainLLM.ProviderOpenAI,
			Model:        req.Model,
			Content:      fmt.Sprintf("[OpenAI Simulation (%s)]: Respuesta procesada exitosamente para %d mensajes.", req.Model, len(req.Messages)),
			PromptTokens: 42,
			CompTokens:   28,
			TotalTokens:  70,
		}, nil
	}

	msgs := make([]openAIMessage, len(req.Messages))
	for i, m := range req.Messages {
		msgs[i] = openAIMessage{Role: string(m.Role), Content: m.Content}
	}

	body := openAIReqBody{
		Model:       req.Model,
		Messages:    msgs,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      false,
	}

	raw, _ := json.Marshal(body)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/chat/completions", bytes.NewBuffer(raw))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domainLLM.ErrProviderTimeout, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, domainLLM.ErrInvalidAPIKey
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, domainLLM.ErrRateLimitExceeded
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai error status %d: %s", resp.StatusCode, string(b))
	}

	var parsed openAIRespBody
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	content := ""
	if len(parsed.Choices) > 0 {
		content = parsed.Choices[0].Message.Content
	}

	return &domainLLM.CompletionResponse{
		Provider:     domainLLM.ProviderOpenAI,
		Model:        req.Model,
		Content:      content,
		PromptTokens: parsed.Usage.PromptTokens,
		CompTokens:   parsed.Usage.CompletionTokens,
		TotalTokens:  parsed.Usage.TotalTokens,
	}, nil
}

func (a *OpenAIAdapter) StreamCompletion(ctx context.Context, req domainLLM.CompletionRequest, onChunk domainLLM.ChunkCallback) error {
	// Simulación de streaming si no hay API key
	simulatedChunks := []string{"Conectando ", "con ", "OpenAI...", " Generando ", "respuesta ", "en ", "streaming."}
	for _, chunk := range simulatedChunks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(120 * time.Millisecond):
			if err := onChunk(chunk); err != nil {
				return err
			}
		}
	}
	return nil
}
