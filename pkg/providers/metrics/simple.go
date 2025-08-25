package metrics

import (
	"sync"
	"time"
)

// MetricsCollector interface - defined locally to avoid import cycle
type MetricsCollector interface {
	IncrementCounter(name string, labels map[string]string)
	RecordHistogram(name string, value float64, labels map[string]string)
	RecordGauge(name string, value float64, labels map[string]string)
	StartTimer(name string) Timer
}

// Timer interface for measuring durations
type Timer interface {
	Stop(labels ...map[string]string)
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Type        string `yaml:"type"`
	Enabled     bool   `yaml:"enabled"`
	Port        int    `yaml:"port"`
	Path        string `yaml:"path"`
	ServiceName string `yaml:"service_name"`
}

// simpleMetrics is a basic in-memory metrics collector
type simpleMetrics struct {
	counters   map[string]int64
	histograms map[string][]float64
	gauges     map[string]float64
	mutex      sync.RWMutex
}

// simpleTimer implements Timer interface
type simpleTimer struct {
	name      string
	startTime time.Time
	metrics   *simpleMetrics
}

// NewSimple creates a new simple metrics collector
func NewSimple(config MetricsConfig) (MetricsCollector, error) {
	return &simpleMetrics{
		counters:   make(map[string]int64),
		histograms: make(map[string][]float64),
		gauges:     make(map[string]float64),
	}, nil
}

// IncrementCounter increments a counter
func (m *simpleMetrics) IncrementCounter(name string, labels map[string]string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.buildKey(name, labels)
	m.counters[key]++
}

// RecordHistogram records a value in a histogram
func (m *simpleMetrics) RecordHistogram(name string, value float64, labels map[string]string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.buildKey(name, labels)
	m.histograms[key] = append(m.histograms[key], value)
}

// RecordGauge sets a gauge value
func (m *simpleMetrics) RecordGauge(name string, value float64, labels map[string]string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.buildKey(name, labels)
	m.gauges[key] = value
}

// StartTimer starts a timer
func (m *simpleMetrics) StartTimer(name string) Timer {
	return &simpleTimer{
		name:      name,
		startTime: time.Now(),
		metrics:   m,
	}
}

// Stop stops the timer and records the duration
func (t *simpleTimer) Stop(labels ...map[string]string) {
	duration := time.Since(t.startTime).Seconds()
	
	var labelMap map[string]string
	if len(labels) > 0 {
		labelMap = labels[0]
	}
	
	t.metrics.RecordHistogram(t.name+"_duration_seconds", duration, labelMap)
}

// buildKey builds a metric key from name and labels
func (m *simpleMetrics) buildKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	
	key := name
	for k, v := range labels {
		key += "," + k + "=" + v
	}
	return key
}

// GetCounters returns all counter values (useful for testing)
func (m *simpleMetrics) GetCounters() map[string]int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	result := make(map[string]int64)
	for k, v := range m.counters {
		result[k] = v
	}
	return result
}

// GetHistograms returns all histogram values (useful for testing)
func (m *simpleMetrics) GetHistograms() map[string][]float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	result := make(map[string][]float64)
	for k, v := range m.histograms {
		values := make([]float64, len(v))
		copy(values, v)
		result[k] = values
	}
	return result
}

// GetGauges returns all gauge values (useful for testing)
func (m *simpleMetrics) GetGauges() map[string]float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	result := make(map[string]float64)
	for k, v := range m.gauges {
		result[k] = v
	}
	return result
}