package config_test

import (
	"reflect"
	"testing"

	config "github.com/cymiam/metrics-store/internal/config/server"
	"github.com/stretchr/testify/require"
)

func TestParseServerConfig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want *config.ServerConfig
	}{
		{
			name: "Deafult Config",
			want: &config.ServerConfig{
				Addr:             "localhost:8080",
				StoreInterval:    300,
				FileStoragePath:  "metrics.json",
				Restore:          false,
				ConnectionString: "postgres://postgres:mysecretpassword@localhost:5432/metrics?sslmode=disable",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := config.ParseServerConfig()
			require.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
