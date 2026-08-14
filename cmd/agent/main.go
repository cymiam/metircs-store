package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/cymiam/metircs-store/internal/agent"
	"github.com/cymiam/metircs-store/internal/logger"
	models "github.com/cymiam/metircs-store/internal/model"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
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

		if time.Since(lastReport) >= time.Duration(agent.Config.ReportInterval*int64(time.Second)) {
			for _, m := range metrics {
				sendMetric(agent.Client, agent.Config.Addr, m)
			}
			randValue := rand.Float64()
			sendMetric(agent.Client, agent.Config.Addr, models.Metrics{ID: "PollCount", MType: "counter", Delta: &agent.PollCount})
			sendMetric(agent.Client, agent.Config.Addr, models.Metrics{ID: "RandomValue", MType: "gauge", Value: &randValue})
			lastReport = time.Now()
		}
		time.Sleep(time.Duration(agent.Config.PollInterval) * time.Second)
	}

}

func sendMetric(client resty.Client, addr string, m models.Metrics) {

	metric, err := easyjson.Marshal(m)
	if err != nil {
		logger.Log.Error("Cannot marshal metric",
			zap.String("Metric Name", m.ID),
		)
		return
	}

	var buffer bytes.Buffer

	w := gzip.NewWriter(&buffer)

	_, err = w.Write(metric)
	if err != nil {
		logger.Log.Error("Error write metric data to buffer", zap.String("Metric Name", m.ID), zap.Error(err))
		return
	}

	err = w.Close()

	if err != nil {
		logger.Log.Error("Error compress metric", zap.String("Metric Name", m.ID), zap.Error(err))
		return
	}

	req := client.R()
	req.Method = "POST"
	req.URL = fmt.Sprintf("http://%s/update/", addr)
	req.Body = buffer.Bytes()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err1 := req.Send()
	if err1 != nil {
		logger.Log.Error("Error sending metric", zap.String("Metric Name", m.ID), zap.Error(err))
		return
	}

	logger.Log.Info("Send metric",
		zap.String("Metric Name", m.ID),
		zap.String("Metric Type", m.MType),
		zap.String("Host", req.URL),
		zap.Int("StatusCode", resp.StatusCode()),
		zap.String("Body", string(resp.Body())))
}
