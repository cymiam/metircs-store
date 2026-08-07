package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Addr string
}

func parseServerFlags(serverConfig *ServerConfig) {
	flag.StringVar(&serverConfig.Addr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}

func ParseServerConfig() ServerConfig {
	var serverConfig ServerConfig
	parseServerFlags(&serverConfig)

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverConfig.Addr = envRunAddr
	}

	return serverConfig
}
