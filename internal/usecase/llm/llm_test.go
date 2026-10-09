package llm_test

import (
	"context"
	"testing"

	domainLLM "proyecto_ia/internal/domain/llm"
	infraLLM "proyecto_ia/internal/infrastructure/llm"
	usecaseLLM "proyecto_ia/internal/usecase/llm"
)

func TestGenerateCompletionUseCase_Decoupled(t *testing.T) {
	// 1. Configurar Factory
	factory := infraLLM.NewProviderFactory()

	// Registrar adaptadores con dummy keys (simulación limpia)
	factory.Register(infraLLM.NewOpenAIAdapter(""))
	factory.Register(infraLLM.NewGeminiAdapter(""))

	// 2. Instanciar Caso de Uso
	uc := usecaseLLM.NewGenerateCompletionUseCase(factory)

	// 3. Probar llamada con OpenAI
	ctx := context.Background()
	reqOpenAI := domainLLM.CompletionRequest{
		Provider: domainLLM.ProviderOpenAI,
		Model:    "gpt-4o",
		Messages: []domainLLM.Message{
			{Role: domainLLM.RoleUser, Content: "Hola arquitectura limpia"},
		},
	}

	respOpenAI, err := uc.Execute(ctx, reqOpenAI)
	if err != nil {
		t.Fatalf("OpenAI execution failed: %v", err)
	}
	if respOpenAI.Provider != domainLLM.ProviderOpenAI {
		t.Errorf("Expected provider %s, got %s", domainLLM.ProviderOpenAI, respOpenAI.Provider)
	}

	// 4. Probar llamada intercambiable con Gemini sin alterar el caso de uso
	reqGemini := domainLLM.CompletionRequest{
		Provider: domainLLM.ProviderGemini,
		Model:    "gemini-1.5-pro",
		Messages: []domainLLM.Message{
			{Role: domainLLM.RoleUser, Content: "Hola Gemini"},
		},
	}

	respGemini, err := uc.Execute(ctx, reqGemini)
	if err != nil {
		t.Fatalf("Gemini execution failed: %v", err)
	}
	if respGemini.Provider != domainLLM.ProviderGemini {
		t.Errorf("Expected provider %s, got %s", domainLLM.ProviderGemini, respGemini.Provider)
	}
}
