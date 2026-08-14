package config_test

import (
	"reflect"
	"testing"

	config "github.com/cymiam/metircs-store/internal/config/server"
)

func TestParseServerConfig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want config.ServerConfig
	}{
		{
			name: "Deafult Config",
			want: config.ServerConfig{
				Addr:            "localhost:8080",
				StoreInterval:   300,
				FileStoragePath: "metrics.json",
				Restore:         false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: update the condition below to compare got with tt.want.
			if got := config.ParseServerConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
