package wails

import (
	"context"
	"fmt"

	domainLLM "proyecto_ia/internal/domain/llm"
	usecaseLLM "proyecto_ia/internal/usecase/llm"
)

// LLMHandler expone casos de uso LLM al frontend mediante Wails Bindings
type LLMHandler struct {
	genUseCase    *usecaseLLM.GenerateCompletionUseCase
	streamUseCase *usecaseLLM.StreamCompletionUseCase
	app           *App
}

func NewLLMHandler(
	genUseCase *usecaseLLM.GenerateCompletionUseCase,
	streamUseCase *usecaseLLM.StreamCompletionUseCase,
	app *App,
) *LLMHandler {
	return &LLMHandler{
		genUseCase:    genUseCase,
		streamUseCase: streamUseCase,
		app:           app,
	}
}

// GenerateCompletion método invocable directamente desde JavaScript: window.go.wails.LLMHandler.GenerateCompletion(...)
func (h *LLMHandler) GenerateCompletion(req domainLLM.CompletionRequest) (*domainLLM.CompletionResponse, error) {
	ctx := context.Background()
	if h.app.ctx != nil {
		ctx = h.app.ctx
	}

	return h.genUseCase.Execute(ctx, req)
}

// StreamPrompt inicia el streaming emitiendo eventos reactivos 'llm:stream:chunk' hacia el frontend
func (h *LLMHandler) StreamPrompt(sessionID string, req domainLLM.CompletionRequest) error {
	ctx := context.Background()
	if h.app.ctx != nil {
		ctx = h.app.ctx
	}

	return h.streamUseCase.Execute(ctx, req, func(chunk string) error {
		h.app.Emit(fmt.Sprintf("llm:stream:%s", sessionID), chunk)
		return nil
	})
}
