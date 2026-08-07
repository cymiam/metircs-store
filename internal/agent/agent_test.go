package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgent_PollRuntimeMetrics(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want int64
	}{
		{
			name: "Test metric count",
			want: 26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAgent()
			got := a.PollRuntimeMetrics()
			// TODO: update the condition below to compare got with tt.want.
			if len(got) != int(tt.want) {
				t.Errorf("PollRuntimeMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgent_PollRuntimeMetricsPresent(t *testing.T) {
	tests := []struct {
		name  string
		names []string // description of this test case
	}{
		{
			name: "Metrics Present",
			names: []string{"Alloc",
				"BuckHashSys",
				"Frees",
				"GCCPUFraction",
				"GCSys",
				"HeapAlloc",
				"HeapIdle",
				"HeapInuse",
				"HeapObjects",
				"HeapReleased",
				"HeapSys",
				"LastGC",
				"Lookups",
				"MCacheInuse",
				"MCacheSys",
				"MSpanSys",
				"Mallocs",
				"NextGC",
				"NumForcedGC",
				"NumGC",
				"OtherSys",
				"PauseTotalNs",
				"StackInuse",
				"StackSys",
				"Sys",
				"TotalAlloc"},
		},
	}

	a := NewAgent()
	metrics := a.PollRuntimeMetrics()

	got := make([]string, 0)

	for _, metric := range metrics {
		got = append(got, metric.Name)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, name := range tt.names {
				assert.Contains(t, got, name)
			}
		})
	}
}

func TestAgent_UpdatePollCount(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want int64
	}{
		{
			name: "Test Update Pollcount",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAgent()
			a.PollRuntimeMetrics()
			got := a.PollCount
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("PollRuntimeMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
