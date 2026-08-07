package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type AgentConfig struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

func parseAgentFlags(agentConfig *AgentConfig) {
	flag.StringVar(&agentConfig.Addr, "a", "localhost:8080", "address and port of server")
	flag.Int64Var(&agentConfig.ReportInterval, "r", 10, "report interval")
	flag.Int64Var(&agentConfig.PollInterval, "p", 2, "poll interval")
	flag.Parse()
}

func ParseAgentConfig() AgentConfig {
	var agentConfig AgentConfig
	parseAgentFlags(&agentConfig)
	err := env.Parse(&agentConfig)

	if err != nil {
		log.Fatal(err)
	}

	return agentConfig
}
