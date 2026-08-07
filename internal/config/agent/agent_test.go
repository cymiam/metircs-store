package config

import (
	"reflect"
	"testing"
)

func TestParseAgentConfig(t *testing.T) {
	tests := []struct {
		name string
		want AgentConfig
	}{
		{
			name: "Deafult Config",
			want: AgentConfig{
				Addr:           "localhost:8080",
				ReportInterval: 10,
				PollInterval:   2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseAgentConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAgentConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
