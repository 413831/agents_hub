package llm

import (
	"fmt"
	"sync"

	domainLLM "proyecto_ia/internal/domain/llm"
)

// ProviderFactory gestiona el registro y resolución de adaptadores LLM
type ProviderFactory struct {
	mu        sync.RWMutex
	providers map[domainLLM.ProviderType]domainLLM.LLMProvider
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[domainLLM.ProviderType]domainLLM.LLMProvider),
	}
}

func (f *ProviderFactory) Register(provider domainLLM.LLMProvider) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.providers[provider.Type()] = provider
}

func (f *ProviderFactory) Get(pType domainLLM.ProviderType) (domainLLM.LLMProvider, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	provider, exists := f.providers[pType]
	if !exists {
		return nil, fmt.Errorf("%w: '%s'", domainLLM.ErrProviderNotFound, pType)
	}
	return provider, nil
}

func (f *ProviderFactory) ListAvailable() []domainLLM.ProviderType {
	f.mu.RLock()
	defer f.mu.RUnlock()

	list := make([]domainLLM.ProviderType, 0, len(f.providers))
	for p := range f.providers {
		list = append(list, p)
	}
	return list
}
