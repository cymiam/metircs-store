package agent

import (
	"testing"
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
