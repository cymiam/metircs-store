package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/cymiam/metircs-store/internal/agent"
)

func main() {
	client := http.Client{
		Timeout: time.Second * 1, // интервал ожидания: 1 секунда
	}
	agent := agent.NewAgent()

	pollCounter := 0
	for {
		metrics := agent.PollRuntimeMetrics()
		pollCounter++

		if pollCounter == 5 {
			for _, metric := range metrics {
				sendMetric(&client, "gauge", metric.Name, fmt.Sprintf("%f", metric.Value))
			}
			sendMetric(&client, "counter", "PollCount", fmt.Sprintf("%d", agent.PollCount))
			sendMetric(&client, "gauge", "RandomValue", fmt.Sprintf("%f", rand.Float64()))
			pollCounter = 0
		}
		time.Sleep(2 * time.Second)
	}

}

func sendMetric(client *http.Client, metricType, metricName, metricValue string) {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, metricName, metricValue)

	_, err := client.Post(url, "Content-Type: text/plain", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Отправил метрику %s, со значением %s\n", metricName, metricValue)
}
