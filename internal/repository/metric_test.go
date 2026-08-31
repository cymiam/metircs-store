package repository

import (
	"reflect"
	"testing"

	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "Test New Mem storage",
			want: &MemStorage{
				Gauges:   make(map[string]float64),
				Counters: make(map[string]int64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetMetric(t *testing.T) {

	m := &MemStorage{
		Counters: map[string]int64{
			"testCounter": int64(5),
		},
		Gauges: map[string]float64{
			"testGauge": 3.14,
		},
	}

	counterValue := int64(5)
	gaugeValue := 3.14

	tests := []struct {
		name       string
		metricName string
		metricType string
		want       models.Metric
		wantErr    bool
	}{
		{
			name:       "Get existing counter",
			metricName: "testCounter",
			metricType: "counter",
			want:       models.Metric{ID: "testCounter", MType: "counter", Delta: &counterValue},
			wantErr:    false,
		},
		{
			name:       "Get non existing counter",
			metricName: "nonCounter",
			metricType: "counter",
			want:       models.Metric{},
			wantErr:    true,
		},
		{
			name:       "Get existing gauge",
			metricName: "testGauge",
			metricType: "gauge",
			want:       models.Metric{ID: "testGauge", MType: "gauge", Value: &gaugeValue},
			wantErr:    false,
		},
		{
			name:       "Get non existing gauge",
			metricName: "nonGauge",
			metricType: "gauge",
			want:       models.Metric{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.GetMetric(t.Context(), tt.metricName, tt.metricType)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	m := &MemStorage{
		Counters: map[string]int64{
			"testCounter": int64(5),
		},
		Gauges: map[string]float64{
			"testGauge": 3.14,
		},
	}

	counterValue := int64(5)
	gaugeValue := 3.14

	tests := []struct {
		name       string
		metricName string
		metricType string
		want       []models.Metric
	}{
		{
			name:       "Get all metrics",
			metricName: "testCounter",
			metricType: "counter",
			want: []models.Metric{
				{ID: "testCounter", MType: "counter", Delta: &counterValue},
				{ID: "testGauge", MType: "gauge", Value: &gaugeValue},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.GetAll(t.Context())

			require.NoError(t, err)

			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestMemStorage_SetMetric(t *testing.T) {

	counterValue := int64(5)
	gaugeValue := 3.14

	tests := []struct {
		name    string
		metric  models.Metric
		wantErr bool
	}{
		{
			name: "Create counter",
			metric: models.Metric{
				ID:    "TestCounter",
				MType: "counter",
				Delta: &counterValue,
			},
			wantErr: false,
		},
		{
			name: "Create counter no delta",
			metric: models.Metric{
				ID:    "TestCounter2",
				MType: "counter",
			},
			wantErr: true,
		},
		{
			name: "Create Gauge",
			metric: models.Metric{
				ID:    "TestGauge",
				MType: "gauge",
				Value: &gaugeValue,
			},
			wantErr: false,
		},
		{
			name: "Create gauge no value",
			metric: models.Metric{
				ID:    "TestGauge2",
				MType: "gauge",
			},
			wantErr: true,
		},
		{
			name: "Create Unknown type",
			metric: models.Metric{
				ID:    "TestUnknown",
				MType: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewStore()
			err := m.SetMetric(t.Context(), tt.metric)

			if tt.wantErr {
				require.Error(t, err)

				metrics, err := m.GetAll(t.Context())
				require.NoError(t, err)
				require.Empty(t, metrics)

				return
			}

			require.NoError(t, err)

			got, err := m.GetMetric(
				t.Context(),
				tt.metric.ID,
				tt.metric.MType,
			)
			require.NoError(t, err)
			require.Equal(t, tt.metric, got)
		})
	}
}
