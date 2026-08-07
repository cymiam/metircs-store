package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/cymiam/metircs-store/internal/agent"
	"github.com/cymiam/metircs-store/internal/config"
)

func main() {
	agent := agent.NewAgent()
	config.ParseAgentFlags()
	lastReport := time.Now()
	for {
		metrics := agent.PollRuntimeMetrics()

		if time.Since(lastReport) >= time.Duration(config.AgentConfig.ReportInterval) {
			for _, metric := range metrics {
				sendMetric(&agent.Client, "gauge", metric.Name, fmt.Sprintf("%f", metric.Value), config.AgentConfig.Addr)
			}
			sendMetric(&agent.Client, "counter", "PollCount", fmt.Sprintf("%d", agent.PollCount), config.AgentConfig.Addr)
			sendMetric(&agent.Client, "gauge", "RandomValue", fmt.Sprintf("%f", rand.Float64()), config.AgentConfig.Addr)
			lastReport = time.Now()
		}
		time.Sleep(time.Duration(config.AgentConfig.PollInterval) * time.Second)
	}

}

func sendMetric(client *http.Client, metricType, metricName, metricValue, host string) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", host, metricType, metricName, metricValue)

	resp, err := client.Post(url, "Content-Type: text/plain", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Отправил метрику %s, со значением %s\n", metricName, metricValue)
}
