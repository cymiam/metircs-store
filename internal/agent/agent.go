package agent

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync/atomic"
	"time"

	config "github.com/cymiam/metrics-store/internal/config/agent"
	"github.com/cymiam/metrics-store/internal/errors/agenterrors"
	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/pkg/compress"
	"github.com/cymiam/metrics-store/pkg/hmac"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

type Agent struct {
	PollCount int64
	Client    resty.Client
	Config    config.AgentConfig
	Logger    *zap.Logger
}

func NewAgent(logger *zap.Logger, cfg config.AgentConfig) *Agent {

	return &Agent{
		PollCount: 0,
		Client:    *resty.New(),
		Config:    cfg,
		Logger:    logger,
	}
}

func (a *Agent) PollRuntimeMetrics() models.Metrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := []models.Metric{models.Metric{ID: "Alloc", Value: helper(float64(m.Alloc)), MType: "gauge"}}
	metrics = append(metrics, models.Metric{ID: "BuckHashSys", Value: helper(float64(m.BuckHashSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "Frees", Value: helper(float64(m.Frees)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "GCCPUFraction", Value: helper(float64(m.GCCPUFraction)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "GCSys", Value: helper(float64(m.GCSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapAlloc", Value: helper(float64(m.HeapAlloc)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapIdle", Value: helper(float64(m.HeapIdle)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapInuse", Value: helper(float64(m.HeapInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapObjects", Value: helper(float64(m.HeapObjects)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapReleased", Value: helper(float64(m.HeapReleased)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapSys", Value: helper(float64(m.HeapSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "LastGC", Value: helper(float64(m.LastGC)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "Lookups", Value: helper(float64(m.Lookups)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "MCacheInuse", Value: helper(float64(m.MCacheInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "MCacheSys", Value: helper(float64(m.MCacheSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "MSpanInuse", Value: helper(float64(m.MSpanInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "MSpanSys", Value: helper(float64(m.MSpanSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "Mallocs", Value: helper(float64(m.Mallocs)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "NextGC", Value: helper(float64(m.NextGC)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "NumForcedGC", Value: helper(float64(m.NumForcedGC)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "NumGC", Value: helper(float64(m.NumGC)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "OtherSys", Value: helper(float64(m.OtherSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "PauseTotalNs", Value: helper(float64(m.PauseTotalNs)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "StackInuse", Value: helper(float64(m.StackInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "StackSys", Value: helper(float64(m.StackSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "Sys", Value: helper(float64(m.Sys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "TotalAlloc", Value: helper(float64(m.TotalAlloc)), MType: "gauge"})

	a.PollCount++
	return metrics
}

func helper(f float64) *float64 {
	return &f
}

// Collect Runtime metrics into channel
func (a *Agent) CollectRuntimeMetrics(ctx context.Context, metricsChan chan<- models.Metric) error {
	pollInterval := time.Duration(a.Config.PollInterval) * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			metrics := a.PollRuntimeMetrics()
			for _, m := range metrics {
				select {
				case metricsChan <- m:
				case <-ctx.Done():
					return nil
				}
			}
			newCount := atomic.AddInt64(&a.PollCount, 1)

			select {
			case metricsChan <- models.Metric{ID: "PollCount", MType: "counter", Delta: &newCount}:
			case <-ctx.Done():
				return nil
			}

			randVal := rand.Float64()
			select {
			case metricsChan <- models.Metric{ID: "RandomValue", MType: "gauge", Value: &randVal}:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

// Collect TotalMemory FreeMemory CPUutil
func (a *Agent) CollectUtilizationMetrics(ctx context.Context, metricsChan chan<- models.Metric) error {
	pollInterval := time.Duration(a.Config.PollInterval) * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:

			memStats, _ := mem.VirtualMemory()
			loadStats, _ := load.Avg()

			totalMemory := float64(memStats.Total)
			freeMemory := float64(memStats.Available)
			cpuUtil := loadStats.Load1

			select {
			case metricsChan <- models.Metric{ID: "TotalMemory", MType: "gauge", Value: &totalMemory}:
			case <-ctx.Done():
				return nil
			}

			select {
			case metricsChan <- models.Metric{ID: "FreeMemory", MType: "gauge", Value: &freeMemory}:
			case <-ctx.Done():
				return nil
			}

			select {
			case metricsChan <- models.Metric{ID: "CPUutilization1", MType: "gauge", Value: &cpuUtil}:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func sendMetrics(client *resty.Client, addr string, m models.Metrics, logger *zap.Logger, key string) error {
	metrics, err := easyjson.Marshal(m)
	if err != nil {
		return fmt.Errorf("cannot marshal metric batch: %w", err)
	}

	gzipped, err := compress.GzipCompress(metrics)
	if err != nil {
		return fmt.Errorf("cannot compress metric: %w", err)
	}

	req := client.R()
	req.Method = "POST"
	req.URL = fmt.Sprintf("http://%s/updates/", addr)
	req.SetBody(gzipped)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	if key != "" {
		req.Header.Set("HashSHA256", hmac.CalculateSha256Sum(gzipped, key))
	}

	resp, err := req.Send()
	if err != nil {
		// Классифицируем ошибку сети
		return agenterrors.ClassifyAgentError(err)
	}

	// Проверяем HTTP-статус ответа (например, 500 тоже стоит считать ошибкой)
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("server returned error status: %d", resp.StatusCode())
	}

	logger.Info("Successfully sent metrics",
		zap.Int("metric_count", len(m)),
		zap.Int("server_response", resp.StatusCode()),
		zap.String("hash", req.Header.Get("HashSHA256")),
	)

	return nil
}

func SendMetricsWithRetry(client *resty.Client, addr string, m models.Metrics, logger *zap.Logger, key string) error {
	var lastErr error

	intervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for attempt := 0; attempt < len(intervals); attempt++ {
		err := sendMetrics(client, addr, m, logger, key)
		if err == nil {
			return nil // Успех, выходим из функции
		}

		lastErr = err
		classification := agenterrors.Classify(err)

		if classification == agenterrors.NonRetriable {
			logger.Error("Non-retriable error while sending metrics", zap.Error(err))
			return err // Прерываем ретраи сразу
		}

		// Если это не последняя попытка, ждем перед следующей
		if attempt < len(intervals)-1 {
			logger.Info("Retrying to send metrics", zap.Int("attempt", attempt+1), zap.Error(err))
			time.Sleep(intervals[attempt])
		}
	}

	logger.Error("Could not send metrics, all retries exhausted", zap.Error(lastErr))
	return lastErr

}

func (a *Agent) StartAggregator(
	ctx context.Context,
	metricsIn <-chan models.Metric,
	jobsOut chan<- models.Metrics,
) error {

	reportInterval := time.Duration(a.Config.ReportInterval) * time.Second
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()

	var batch models.Metrics

	for {
		select {
		case <-ctx.Done():
			close(jobsOut) // Сигнализируем воркерам, что новых задач не будет
			return nil

		case m := <-metricsIn:
			batch = append(batch, m)

		case <-ticker.C:
			if len(batch) > 0 {
				select {
				case jobsOut <- batch:
					batch = nil
				case <-ctx.Done():
					return nil
				}
			}
		}
	}
}

func (a *Agent) StartWorkerPool(ctx context.Context, jobsIn <-chan models.Metrics, numOfWorkers int) error {
	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < numOfWorkers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case batch, ok := <-jobsIn:
					if !ok {
						return nil
					}
					err := SendMetricsWithRetry(&a.Client, a.Config.Addr, batch, a.Logger, a.Config.Key)

					if err != nil {
						a.Logger.Error("cannot send request", zap.Int("WorkerID", i), zap.Error(err))
						return err
					}
				}
			}
		})
	}

	return g.Wait()
}
