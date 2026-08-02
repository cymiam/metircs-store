package config

import "flag"

type serverConfig struct {
	Addr string
}

var ServerConfig serverConfig

func ParseServerFlags() {
	flag.StringVar(&ServerConfig.Addr, "a", ":8080", "address and port to run server")
	flag.Parse()
}
