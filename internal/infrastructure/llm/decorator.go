package llm

import (
	"context"
	"time"

	domainLLM "proyecto_ia/internal/domain/llm"
)

// MetricsRecorder interfaz para desacoplar el recolector de métricas
type MetricsRecorder interface {
	RecordLLMRequest(provider string, model string, duration time.Duration, tokens int, success bool)
}

// TelemetryDecorator decora cualquier LLMProvider añadiendo métricas de observabilidad
type TelemetryDecorator struct {
	inner    domainLLM.LLMProvider
	recorder MetricsRecorder
}

func NewTelemetryDecorator(inner domainLLM.LLMProvider, recorder MetricsRecorder) *TelemetryDecorator {
	return &TelemetryDecorator{
		inner:    inner,
		recorder: recorder,
	}
}

func (d *TelemetryDecorator) Type() domainLLM.ProviderType {
	return d.inner.Type()
}

func (d *TelemetryDecorator) SupportedModels() []string {
	return d.inner.SupportedModels()
}

func (d *TelemetryDecorator) GenerateCompletion(ctx context.Context, req domainLLM.CompletionRequest) (*domainLLM.CompletionResponse, error) {
	start := time.Now()
	res, err := d.inner.GenerateCompletion(ctx, req)
	elapsed := time.Since(start)

	if d.recorder != nil {
		tokens := 0
		if res != nil {
			tokens = res.TotalTokens
		}
		d.recorder.RecordLLMRequest(string(req.Provider), req.Model, elapsed, tokens, err == nil)
	}

	return res, err
}

func (d *TelemetryDecorator) StreamCompletion(ctx context.Context, req domainLLM.CompletionRequest, onChunk domainLLM.ChunkCallback) error {
	start := time.Now()
	err := d.inner.StreamCompletion(ctx, req, onChunk)
	elapsed := time.Since(start)

	if d.recorder != nil {
		d.recorder.RecordLLMRequest(string(req.Provider), req.Model, elapsed, 0, err == nil)
	}

	return err
}
