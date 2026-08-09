package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/cymiam/metircs-store/internal/agent"
	"github.com/cymiam/metircs-store/internal/logger"
	"go.uber.org/zap"
)

func main() {

	if err := logger.Initialize("Info"); err != nil {
		panic(err)
	}
	agent := agent.NewAgent()
	lastReport := time.Now()
	for {
		metrics := agent.PollRuntimeMetrics()

		if time.Since(lastReport) >= time.Duration(agent.Config.ReportInterval) {
			for _, metric := range metrics {
				sendMetric(&agent.Client, "gauge", metric.Name, fmt.Sprintf("%f", metric.Value), agent.Config.Addr)
			}
			sendMetric(&agent.Client, "counter", "PollCount", fmt.Sprintf("%d", agent.PollCount), agent.Config.Addr)
			sendMetric(&agent.Client, "gauge", "RandomValue", fmt.Sprintf("%f", rand.Float64()), agent.Config.Addr)
			lastReport = time.Now()
		}
		time.Sleep(time.Duration(agent.Config.PollInterval) * time.Second)
	}

}

func sendMetric(client *http.Client, metricType, metricName, metricValue, host string) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", host, metricType, metricName, metricValue)

	resp, err := client.Post(url, "Content-Type: text/plain", nil)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer resp.Body.Close()
	logger.Log.Info("Отправил метрику",
		zap.String("Metric Name", metricName),
		zap.String("Metric Value", metricValue))
}
