package main

import (
	"log"
	"net/http"

	"github.com/cymiam/metircs-store/internal/handler"
)

func main() {
	r := handler.MetricRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
