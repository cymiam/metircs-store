package config

import (
	"flag"
)

type agentConfig struct {
	Addr           string
	ReportInterval int64
	PollInterval   int64
}

var AgentConfig agentConfig

func ParseAgentFlags() {
	flag.StringVar(&AgentConfig.Addr, "a", ":8080", "address and port of server")
	flag.Int64Var(&AgentConfig.PollInterval, "r", 10, "report interval")
	flag.Int64Var(&AgentConfig.PollInterval, "p", 10, "poll interval")
	flag.Parse()
}
