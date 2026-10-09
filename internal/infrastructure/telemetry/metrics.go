package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// MetricsSnapshot provee una vista estructurada de métricas para la UI y Grafana
type MetricsSnapshot struct {
	TotalRequests   int64              `json:"total_requests"`
	TotalTokens     int64              `json:"total_tokens"`
	SuccessRate     float64            `json:"success_rate"`
	AverageLatency  float64            `json:"average_latency_ms"`
	ActiveWorkers   int32              `json:"active_workers"`
	ProviderCounts  map[string]int64   `json:"provider_counts"`
}

// PrometheusMetricsManager gestiona el registro de métricas y la exposición HTTP
type PrometheusMetricsManager struct {
	mu             sync.RWMutex
	totalRequests  int64
	successCount   int64
	totalTokens    int64
	totalLatencyMs int64
	activeWorkers  int32
	providerCounts map[string]int64
	server         *http.Server
}

func NewPrometheusMetricsManager() *PrometheusMetricsManager {
	return &PrometheusMetricsManager{
		providerCounts: make(map[string]int64),
	}
}

func (m *PrometheusMetricsManager) RecordLLMRequest(provider string, model string, duration time.Duration, tokens int, success bool) {
	atomic.AddInt64(&m.totalRequests, 1)
	if success {
		atomic.AddInt64(&m.successCount, 1)
	}
	atomic.AddInt64(&m.totalTokens, int64(tokens))
	atomic.AddInt64(&m.totalLatencyMs, duration.Milliseconds())

	m.mu.Lock()
	m.providerCounts[provider]++
	m.mu.Unlock()
}

func (m *PrometheusMetricsManager) SetActiveWorkers(count int32) {
	atomic.StoreInt32(&m.activeWorkers, count)
}

func (m *PrometheusMetricsManager) GetSnapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := atomic.LoadInt64(&m.totalRequests)
	success := atomic.LoadInt64(&m.successCount)
	tokens := atomic.LoadInt64(&m.totalTokens)
	latency := atomic.LoadInt64(&m.totalLatencyMs)

	var successRate float64 = 100.0
	if total > 0 {
		successRate = (float64(success) / float64(total)) * 100.0
	}

	var avgLatency float64 = 0.0
	if total > 0 {
		avgLatency = float64(latency) / float64(total)
	}

	countsCopy := make(map[string]int64)
	for k, v := range m.providerCounts {
		countsCopy[k] = v
	}

	return MetricsSnapshot{
		TotalRequests:  total,
		TotalTokens:    tokens,
		SuccessRate:    successRate,
		AverageLatency: avgLatency,
		ActiveWorkers:  atomic.LoadInt32(&m.activeWorkers),
		ProviderCounts: countsCopy,
	}
}

// StartMetricsServer arranca el endpoint /metrics en formato Prometheus estándar para Grafana
func (m *PrometheusMetricsManager) StartMetricsServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		snap := m.GetSnapshot()
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "# HELP llm_requests_total Total de llamadas a modelos LLM\n")
		fmt.Fprintf(w, "# TYPE llm_requests_total counter\n")
		fmt.Fprintf(w, "llm_requests_total %d\n", snap.TotalRequests)

		fmt.Fprintf(w, "# HELP llm_tokens_consumed_total Total de tokens procesados\n")
		fmt.Fprintf(w, "# TYPE llm_tokens_consumed_total counter\n")
		fmt.Fprintf(w, "llm_tokens_consumed_total %d\n", snap.TotalTokens)

		fmt.Fprintf(w, "# HELP llm_request_latency_avg_ms Latencia promedio en milisegundos\n")
		fmt.Fprintf(w, "# TYPE llm_request_latency_avg_ms gauge\n")
		fmt.Fprintf(w, "llm_request_latency_avg_ms %.2f\n", snap.AverageLatency)

		fmt.Fprintf(w, "# HELP background_active_workers Cantidad de workers activos\n")
		fmt.Fprintf(w, "# TYPE background_active_workers gauge\n")
		fmt.Fprintf(w, "background_active_workers %d\n", snap.ActiveWorkers)

		for prov, count := range snap.ProviderCounts {
			fmt.Fprintf(w, "llm_provider_requests_total{provider=\"%s\"} %d\n", prov, count)
		}
	})

	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		_ = m.server.ListenAndServe()
	}()
}

func (m *PrometheusMetricsManager) Stop(ctx context.Context) error {
	if m.server != nil {
		return m.server.Shutdown(ctx)
	}
	return nil
}
