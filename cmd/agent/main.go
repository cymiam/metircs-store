package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/cymiam/metrics-store/internal/agent"
	config "github.com/cymiam/metrics-store/internal/config/agent"
	"github.com/cymiam/metrics-store/internal/errors/agenterrors"
	"github.com/cymiam/metrics-store/internal/logger"
	models "github.com/cymiam/metrics-store/internal/model"
	compress "github.com/cymiam/metrics-store/pkg/compress"
	"github.com/cymiam/metrics-store/pkg/hmac"
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

	retryInterval := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	agent := agent.NewAgent(agentConfig)
	lastReport := time.Now()
	for {
		metrics := agent.PollRuntimeMetrics()

		if time.Since(lastReport) >= time.Duration(agent.Config.ReportInterval*int64(time.Second)) {

			randValue := rand.Float64()
			metrics = append(metrics, models.Metric{ID: "PollCount", MType: "counter", Delta: &agent.PollCount})
			metrics = append(metrics, models.Metric{ID: "RandomValue", MType: "gauge", Value: &randValue})
			for attempt := 0; attempt <= 3; attempt++ {
				err := sendMetrics(agent.Client, agent.Config.Addr, metrics, logger, agentConfig.Key)

				if err == nil {
					lastReport = time.Now()
					break
				}
				classification := agenterrors.Classify(err)

				if classification == agenterrors.NonRetriable {
					logger.Error("couldnt send metrics", zap.Error(err))
					break
				}

				if attempt == len(retryInterval) {
					logger.Error("couldnt send metric, all retries exhausted", zap.Error(err))
					break
				}

				time.Sleep(retryInterval[attempt])

			}
		}
		time.Sleep(time.Duration(agent.Config.PollInterval) * time.Second)
	}

}

func sendMetrics(client resty.Client, addr string, m models.Metrics, logger *zap.Logger, key string) error {

	metrics, err := easyjson.Marshal(m)
	if err != nil {
		return fmt.Errorf("cannot marshal metric batch, %w", err)
	}

	gziped, err := compress.GzipCompress(metrics)

	if err != nil {
		return fmt.Errorf("cannot compress metric, %w", err)
	}

	req := client.R()
	req.Method = "POST"
	req.URL = fmt.Sprintf("http://%s/updates/", addr)
	req.SetBody(gziped)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	if key != "" {
		req.Header.Set("HashSHA256", hmac.CalculateSha256Sum(gziped, key))
	}
	resp, err := req.Send()
	if err != nil {
		err = agenterrors.ClassifyAgentError(err)
		return err
	}

	logger.Info("Sended metrics", zap.Int("metric count", len(m)), zap.Int("server response", resp.StatusCode()), zap.String("hash", req.Header.Get("HashSHA256")))
	return nil
}
