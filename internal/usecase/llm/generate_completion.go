package llm

import (
	"context"
	"fmt"
	"time"

	domainLLM "proyecto_ia/internal/domain/llm"
)

// ProviderResolver interface para desacoplar la resolución de proveedores
type ProviderResolver interface {
	Get(pType domainLLM.ProviderType) (domainLLM.LLMProvider, error)
}

// GenerateCompletionUseCase orquesta la inferencia síncrona
type GenerateCompletionUseCase struct {
	resolver ProviderResolver
}

func NewGenerateCompletionUseCase(resolver ProviderResolver) *GenerateCompletionUseCase {
	return &GenerateCompletionUseCase{
		resolver: resolver,
	}
}

func (uc *GenerateCompletionUseCase) Execute(ctx context.Context, req domainLLM.CompletionRequest) (*domainLLM.CompletionResponse, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("request must contain at least one message")
	}

	provider, err := uc.resolver.Get(req.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve provider %s: %w", req.Provider, err)
	}

	start := time.Now()
	res, err := provider.GenerateCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("provider execution error: %w", err)
	}

	res.DurationMs = time.Since(start).Milliseconds()
	return res, nil
}
