package repository

import (
	"reflect"
	"testing"
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
				Counters: make(map[string][]int64),
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

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		gauges   map[string]float64
		counters map[string][]int64
	}

	tests := []struct {
		name       string
		fields     fields
		metricName string
		want       []int64
		want1      bool
	}{
		{
			name: "Test Positive",
			fields: fields{
				counters: map[string][]int64{
					"test": []int64{1, 2, 3},
				},
				gauges: map[string]float64{
					"test2": 3.14,
				},
			},
			metricName: "test",
			want:       []int64{1, 2, 3},
			want1:      true,
		},
		{
			name: "Test Negative",
			fields: fields{
				counters: map[string][]int64{
					"test": []int64{1, 2, 3},
				},
				gauges: map[string]float64{
					"test2": 3.14,
				},
			},
			metricName: "test2",
			want:       nil,
			want1:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.fields.gauges,
				Counters: tt.fields.counters,
			}
			got, got1 := m.GetCounter(tt.metricName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetCounter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MemStorage.GetCounter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		gauges   map[string]float64
		counters map[string][]int64
	}
	tests := []struct {
		name       string
		fields     fields
		metricName string
		want       float64
		want1      bool
	}{
		{
			name: "Test Positive",
			fields: fields{
				counters: map[string][]int64{
					"test": []int64{1, 2, 3},
				},
				gauges: map[string]float64{
					"test2": 3.14,
				},
			},
			metricName: "test2",
			want:       3.14,
			want1:      true,
		},
		{
			name: "Test Negative",
			fields: fields{
				counters: map[string][]int64{
					"test": []int64{1, 2, 3},
				},
				gauges: map[string]float64{
					"test2": 3.14,
				},
			},
			metricName: "test",
			want:       0.0,
			want1:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.fields.gauges,
				Counters: tt.fields.counters,
			}
			got, got1 := m.GetGauge(tt.metricName)
			if got != tt.want {
				t.Errorf("MemStorage.GetGauge() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MemStorage.GetGauge() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemStorage_GetGauges(t *testing.T) {
	type fields struct {
		gauges   map[string]float64
		counters map[string][]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]float64
	}{
		{
			name: "Test Positive",
			fields: fields{
				counters: map[string][]int64{
					"test": []int64{1, 2, 3},
				},
				gauges: map[string]float64{
					"test2": 3.14,
					"test3": 4.5,
				},
			},
			want: map[string]float64{
				"test2": 3.14,
				"test3": 4.5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.fields.gauges,
				Counters: tt.fields.counters,
			}
			if got := m.GetGauges(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetGauges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounters(t *testing.T) {
	type fields struct {
		gauges   map[string]float64
		counters map[string][]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]int64
	}{
		{
			name: "Test Positive",
			fields: fields{
				counters: map[string][]int64{
					"test1": []int64{1, 2, 3},
					"test2": []int64{3, 2, 1},
				},
				gauges: map[string]float64{
					"test2": 3.14,
					"test3": 4.5,
				},
			},
			want: map[string][]int64{
				"test1": []int64{1, 2, 3},
				"test2": []int64{3, 2, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.fields.gauges,
				Counters: tt.fields.counters,
			}
			if got := m.GetCounters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetCounters() = %v, want %v", got, tt.want)
			}
		})
	}
}
