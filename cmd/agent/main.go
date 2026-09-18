package main

import (
	"context"
	"log"

	"github.com/cymiam/metrics-store/internal/agent"
	config "github.com/cymiam/metrics-store/internal/config/agent"
	"github.com/cymiam/metrics-store/internal/logger"
	models "github.com/cymiam/metrics-store/internal/model"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {

	logger, err := logger.NewLogger("info", "agent-logger")

	if err != nil {
		log.Fatal("Cannot start logger", err)
	}

	agentConfig, err := config.ParseAgentConfig()

	if err != nil {
		log.Fatal("Cannot parse agent config", err)
	}

	agent := agent.NewAgent(logger, agentConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	metricsChan := make(chan models.Metric)

	jobsChan := make(chan models.Metrics, agentConfig.RateLimit)

	defer close(metricsChan)

	// Сбор рантайм метрик
	g.Go(func() error {
		return agent.CollectRuntimeMetrics(ctx, metricsChan)
	})

	// Сбор метрик утилизации
	g.Go(func() error {
		return agent.CollectUtilizationMetrics(ctx, metricsChan)
	})

	// Сбор метрик в батч перед отправкой
	g.Go(func() error {
		return agent.StartAggregator(ctx, metricsChan, jobsChan)
	})

	// Отправка метрик
	g.Go(func() error {
		return agent.StartWorkerPool(ctx, jobsChan, agentConfig.RateLimit)
	})

	if err := g.Wait(); err != nil {
		logger.Error("Agent errors", zap.Error(err))
	}

}
