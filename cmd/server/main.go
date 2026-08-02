package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cymiam/metircs-store/internal/config"
	"github.com/cymiam/metircs-store/internal/handler"
)

func main() {
	config.ParseServerFlags()
	r := handler.MetricRouter()
	fmt.Println("Running server on", config.ServerConfig.Addr)
	log.Fatal(http.ListenAndServe(config.ServerConfig.Addr, r))
}
