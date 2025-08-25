package metrics

import ()

// noopMetrics is a metrics collector that doesn't collect anything
type noopMetrics struct{}

// noopTimer is a timer that doesn't measure anything
type noopTimer struct{}

// NewNoop creates a new no-op metrics collector
func NewNoop(config MetricsConfig) (MetricsCollector, error) {
	return &noopMetrics{}, nil
}

// IncrementCounter does nothing
func (m *noopMetrics) IncrementCounter(name string, labels map[string]string) {}

// RecordHistogram does nothing
func (m *noopMetrics) RecordHistogram(name string, value float64, labels map[string]string) {}

// RecordGauge does nothing
func (m *noopMetrics) RecordGauge(name string, value float64, labels map[string]string) {}

// StartTimer returns a no-op timer
func (m *noopMetrics) StartTimer(name string) Timer {
	return &noopTimer{}
}

// Stop does nothing
func (t *noopTimer) Stop(labels ...map[string]string) {}