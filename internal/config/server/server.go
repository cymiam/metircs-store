package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ServerConfig struct {
	Addr            string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

func parseServerFlags(serverConfig *ServerConfig) {
	flag.StringVar(&serverConfig.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&serverConfig.FileStoragePath, "f", "metrics.json", "local metrics storage path")
	flag.IntVar(&serverConfig.StoreInterval, "i", 300, "Metrics update offset in seconds")
	flag.BoolVar(&serverConfig.Restore, "r", false, "Start server with old metrics")
	flag.Parse()
}

func ParseServerConfig() (*ServerConfig, error) {
	var serverConfig ServerConfig
	parseServerFlags(&serverConfig)

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverConfig.Addr = envRunAddr
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		val, err := strconv.Atoi(envStoreInterval)

		if err != nil {
			return nil, fmt.Errorf("Cannot convert STORE_INTERVAL to INT %w", err)
		}
		serverConfig.StoreInterval = val
	}

	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		serverConfig.FileStoragePath = envStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		val, err := strconv.ParseBool(envRestore)

		if err != nil {
			return nil, fmt.Errorf("Cannot convert RESTORE to BOOL %w", err)
		}
		serverConfig.Restore = val

	}

	return &serverConfig, nil
}
