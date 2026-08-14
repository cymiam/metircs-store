package service

import (
	"testing"

	"github.com/cymiam/metircs-store/internal/repository"
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
			service := NewMetricService(repository.NewStore())
			service.store.SetCounter("test", 1)
			service.UpdateCounter(tt.metricName, tt.newValue)

			expected := []int64{1, 3}
			got, _ := service.GetCounter("test")

			assert.Equal(t, expected, got)
		})
	}
}
