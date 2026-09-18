package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type AgentConfig struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func parseAgentFlags(agentConfig *AgentConfig) {
	flag.StringVar(&agentConfig.Addr, "a", "localhost:8080", "address and port of server")
	flag.Int64Var(&agentConfig.ReportInterval, "r", 10, "report interval")
	flag.Int64Var(&agentConfig.PollInterval, "p", 2, "poll interval")
	flag.StringVar(&agentConfig.Key, "k", "", "encryption key")
	flag.IntVar(&agentConfig.RateLimit, "l", 1, "Number of workers")
	flag.Parse()
}

func ParseAgentConfig() (AgentConfig, error) {
	var agentConfig AgentConfig
	parseAgentFlags(&agentConfig)

	if err := env.Parse(&agentConfig); err != nil {
		return AgentConfig{}, fmt.Errorf("parse server environment: %w", err)
	}

	return agentConfig, nil
}
