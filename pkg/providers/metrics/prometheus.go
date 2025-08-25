package metrics

import (
	"fmt"
)

// NewPrometheus creates a new Prometheus metrics collector
func NewPrometheus(config MetricsConfig) (MetricsCollector, error) {
	return nil, fmt.Errorf("Prometheus provider not implemented yet")
}