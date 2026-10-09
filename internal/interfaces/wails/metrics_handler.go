package wails

import (
	infraTelemetry "proyecto_ia/internal/infrastructure/telemetry"
)

// MetricsHandler expone el snapshot de observabilidad a la vista del dashboard
type MetricsHandler struct {
	metricsMgr *infraTelemetry.PrometheusMetricsManager
}

func NewMetricsHandler(metricsMgr *infraTelemetry.PrometheusMetricsManager) *MetricsHandler {
	return &MetricsHandler{
		metricsMgr: metricsMgr,
	}
}

func (h *MetricsHandler) GetMetricsSnapshot() infraTelemetry.MetricsSnapshot {
	return h.metricsMgr.GetSnapshot()
}
