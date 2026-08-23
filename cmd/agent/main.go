package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/cymiam/metrics-store/internal/agent"
	config "github.com/cymiam/metrics-store/internal/config/agent"
	"github.com/cymiam/metrics-store/internal/logger"
	models "github.com/cymiam/metrics-store/internal/model"
	compress "github.com/cymiam/metrics-store/pkg/compress"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
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
	agent := agent.NewAgent(agentConfig)
	lastReport := time.Now()
	for {
		metrics := agent.PollRuntimeMetrics()

		if time.Since(lastReport) >= time.Duration(agent.Config.ReportInterval*int64(time.Second)) {
			for _, m := range metrics {
				sendMetric(agent.Client, agent.Config.Addr, m, logger)
			}
			randValue := rand.Float64()
			sendMetric(agent.Client, agent.Config.Addr, models.Metrics{ID: "PollCount", MType: "counter", Delta: &agent.PollCount}, logger)
			sendMetric(agent.Client, agent.Config.Addr, models.Metrics{ID: "RandomValue", MType: "gauge", Value: &randValue}, logger)
			lastReport = time.Now()
		}
		time.Sleep(time.Duration(agent.Config.PollInterval) * time.Second)
	}

}

func sendMetric(client resty.Client, addr string, m models.Metrics, logger *zap.Logger) {

	metric, err := easyjson.Marshal(m)
	if err != nil {
		logger.Error("Cannot marshal metric",
			zap.String("Metric Name", m.ID),
		)
		return
	}

	gziped, err := compress.GzipCompress(metric)

	if err != nil {
		logger.Error("Error compress metric", zap.String("Metric Name", m.ID), zap.Error(err))
		return
	}

	req := client.R()
	req.Method = "POST"
	req.URL = fmt.Sprintf("http://%s/update/", addr)
	req.Body = gziped
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := req.Send()
	if err != nil {
		logger.Error("Error sending metric", zap.String("Metric Name", m.ID), zap.Error(err))
		return
	}

	logger.Info("Send metric",
		zap.String("Metric Name", m.ID),
		zap.String("Metric Type", m.MType),
		zap.String("Host", req.URL),
		zap.Int("StatusCode", resp.StatusCode()),
		zap.String("Body", string(resp.Body())))
}
