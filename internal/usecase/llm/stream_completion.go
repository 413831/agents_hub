package llm

import (
	"context"
	"fmt"

	domainLLM "proyecto_ia/internal/domain/llm"
)

// StreamCompletionUseCase orquesta el flujo de streaming token a token
type StreamCompletionUseCase struct {
	resolver ProviderResolver
}

func NewStreamCompletionUseCase(resolver ProviderResolver) *StreamCompletionUseCase {
	return &StreamCompletionUseCase{
		resolver: resolver,
	}
}

func (uc *StreamCompletionUseCase) Execute(ctx context.Context, req domainLLM.CompletionRequest, onChunk domainLLM.ChunkCallback) error {
	provider, err := uc.resolver.Get(req.Provider)
	if err != nil {
		return fmt.Errorf("failed to resolve provider %s: %w", req.Provider, err)
	}

	return provider.StreamCompletion(ctx, req, onChunk)
}
