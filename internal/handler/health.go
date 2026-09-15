package handler

import (
	"net/http"

	"github.com/cymiam/metrics-store/internal/service"
	"go.uber.org/zap"
)

type HealthHandler struct {
	healthService *service.HealthService
	logger        *zap.Logger
}

type HealthHandlerParams struct {
	HealthService *service.HealthService
	Logger        *zap.Logger
}

func NewHealthHandler(params HealthHandlerParams) *HealthHandler {
	return &HealthHandler{
		healthService: params.HealthService,
		logger:        params.Logger,
	}
}

func (handler *HealthHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	err := handler.healthService.PingDB()
	if err != nil {
		handler.logger.Error("Error ping db", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
