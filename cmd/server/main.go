package main

import (
	"log"
	"net/http"

	config "github.com/cymiam/metircs-store/internal/config/server"
	"github.com/cymiam/metircs-store/internal/handler"
	"github.com/cymiam/metircs-store/internal/logger"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Initialize("Info"); err != nil {
		panic(err)
	}
	config := config.ParseServerConfig()

	r := handler.MetricRouter()

	logger.Log.Info("Running server", zap.String("address", config.Addr))
	log.Fatal(http.ListenAndServe(config.Addr, r))
}
