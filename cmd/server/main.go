package main

import (
	"net/http"

	"github.com/cymiam/metircs-store/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	metricHandler := handler.NewMetricHandler()
	mux.HandleFunc("POST /update/{metric_type}/{metric_name}/{metric_value}", metricHandler.HandleUpdate)
	err := http.ListenAndServe("127.0.0.1:8080", mux)

	if err != nil {
		panic(err)
	}
}
