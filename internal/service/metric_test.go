package service

import (
	"testing"

	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestMetricService_UpdateCounter(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		metricName string
		newValue   int64
	}{
		{
			name:       "Test Positive",
			metricName: "test",
			newValue:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewMetricService(MetricServiceParams{Store: repository.NewStore()})
			one := int64(1)
			service.store.SetMetric(t.Context(), models.Metric{ID: "test", MType: "counter", Delta: &one})
			service.UpdateCounter(tt.metricName, tt.newValue)

			val := int64(3)
			expected := models.Metric{ID: "test", MType: "counter", Delta: &val}
			got, err := service.GetMetric("test", "counter")

			assert.NoError(t, err)

			assert.Equal(t, expected, got)
		})
	}
}
