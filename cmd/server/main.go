package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/cymiam/metircs-store/internal/config/server"
	"github.com/cymiam/metircs-store/internal/handler"
)

func main() {

	config := config.ParseServerConfig()

	r := handler.MetricRouter()
	fmt.Println("Running server on", config.Addr)
	log.Fatal(http.ListenAndServe(config.Addr, r))
}
