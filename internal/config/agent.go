package config

import (
	"flag"
	"os"
)

type AgentConfig struct {
	Addr           string
	ReportInterval int64
	PollInterval   int64
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

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		agentConfig.Addr = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		agentConfig.Addr = envReportInterval
	}

	if envPollIntervall := os.Getenv("POLL_INTERVAL"); envPollIntervall != "" {
		agentConfig.Addr = envPollIntervall
	}

	return agentConfig
}
