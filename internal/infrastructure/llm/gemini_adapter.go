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

// GeminiAdapter desacopla la comunicación con Google Gemini
type GeminiAdapter struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewGeminiAdapter(apiKey string) *GeminiAdapter {
	if apiKey == "" {
		apiKey = "dummy-key"
	}
	return &GeminiAdapter{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (g *GeminiAdapter) Type() domainLLM.ProviderType {
	return domainLLM.ProviderGemini
}

func (g *GeminiAdapter) SupportedModels() []string {
	return []string{"gemini-1.5-pro", "gemini-1.5-flash", "gemini-2.0-flash"}
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiReqBody struct {
	Contents []geminiContent `json:"contents"`
}

type geminiRespBody struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

func (g *GeminiAdapter) GenerateCompletion(ctx context.Context, req domainLLM.CompletionRequest) (*domainLLM.CompletionResponse, error) {
	if g.apiKey == "" || g.apiKey == "dummy-key" {
		return &domainLLM.CompletionResponse{
			Provider:     domainLLM.ProviderGemini,
			Model:        req.Model,
			Content:      fmt.Sprintf("[Google Gemini Simulation (%s)]: Inferencia generativa completada para %d mensajes.", req.Model, len(req.Messages)),
			PromptTokens: 35,
			CompTokens:   45,
			TotalTokens:  80,
		}, nil
	}

	contents := make([]geminiContent, len(req.Messages))
	for i, m := range req.Messages {
		role := "user"
		if m.Role == domainLLM.RoleAssistant {
			role = "model"
		}
		contents[i] = geminiContent{
			Role:  role,
			Parts: []geminiPart{{Text: m.Content}},
		}
	}

	body := geminiReqBody{Contents: contents}
	raw, _ := json.Marshal(body)

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", g.baseURL, req.Model, g.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domainLLM.ErrProviderTimeout, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini error status %d: %s", resp.StatusCode, string(b))
	}

	var parsed geminiRespBody
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	text := ""
	if len(parsed.Candidates) > 0 && len(parsed.Candidates[0].Content.Parts) > 0 {
		text = parsed.Candidates[0].Content.Parts[0].Text
	}

	return &domainLLM.CompletionResponse{
		Provider:     domainLLM.ProviderGemini,
		Model:        req.Model,
		Content:      text,
		PromptTokens: parsed.UsageMetadata.PromptTokenCount,
		CompTokens:   parsed.UsageMetadata.CandidatesTokenCount,
		TotalTokens:  parsed.UsageMetadata.TotalTokenCount,
	}, nil
}

func (g *GeminiAdapter) StreamCompletion(ctx context.Context, req domainLLM.CompletionRequest, onChunk domainLLM.ChunkCallback) error {
	simulatedChunks := []string{"Google ", "Gemini ", "procesando ", "stream: ", "multimodalidad ", "y ", "velocidad."}
	for _, chunk := range simulatedChunks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			if err := onChunk(chunk); err != nil {
				return err
			}
		}
	}
	return nil
}
