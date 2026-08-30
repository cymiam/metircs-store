package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	Addr             string `env:"ADDRESS"`
	StoreInterval    int    `env:"STORE_INTERVAL"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	Restore          bool   `env:"RESTORE"`
	ConnectionString string `env:"DATABASE_DSN"`
}

func parseServerFlags(serverConfig *ServerConfig) {
	flag.StringVar(&serverConfig.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&serverConfig.FileStoragePath, "f", "metrics.json", "local metrics storage path")
	flag.IntVar(&serverConfig.StoreInterval, "i", 300, "Metrics update offset in seconds")
	flag.BoolVar(&serverConfig.Restore, "r", false, "Start server with old metrics")
	flag.StringVar(&serverConfig.ConnectionString, "d", "postgres://postgres:mysecretpassword@localhost:5432/metrics?sslmode=disable", "Connection string to database")
	flag.Parse()
}

func ParseServerConfig() (*ServerConfig, error) {
	var serverConfig ServerConfig
	parseServerFlags(&serverConfig)

	if err := env.Parse(&serverConfig); err != nil {
		return nil, fmt.Errorf("parse server environment: %w", err)
	}

	return &serverConfig, nil
}
