package service

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/internal/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMetricSaver_AppendMetricsToJson(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metrics.jsonl")

	const existingRecord = "{\"id\":\"existing\",\"type\":\"gauge\",\"value\":1.5}\n"
	require.NoError(t, os.WriteFile(path, []byte(existingRecord), 0600))

	store := repository.NewStore()
	saver := newTestMetricSaver(t, path, 0, store)
	metricService := NewMetricService(MetricServiceParams{
		Store:  store,
		Saver:  saver,
		Logger: zap.NewNop(),
	})

	metricService.UpdateCounter("requests", 2)
	metricService.UpdateCounter("requests", 3)
	metricService.UpdateGauge("temperature", 10.5)
	metricService.UpdateGauge("temperature", 12)

	events := readMetricEvents(t, path)
	require.Len(t, events, 5)
	require.Empty(t, saver.queue)

	requireMetricGauge(t, events[0], "existing", 1.5)
	requireMetricCounter(t, events[1], "requests", 2)
	requireMetricCounter(t, events[2], "requests", 3)
	requireMetricGauge(t, events[3], "temperature", 10.5)
	requireMetricGauge(t, events[4], "temperature", 12)

	require.Equal(t, int64(5), store.Counters["requests"])
}

func TestMetricSaver_RestoreMetrics(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metrics.jsonl")

	const history = "" +
		"{\"id\":\"requests\",\"type\":\"counter\",\"delta\":2}\n" +
		"{\"id\":\"temperature\",\"type\":\"gauge\",\"value\":10.5}\n" +
		"{\"id\":\"requests\",\"type\":\"counter\",\"delta\":3}\n" +
		"{\"id\":\"temperature\",\"type\":\"gauge\",\"value\":12}\n" +
		"{\"id\":\"requests\",\"type\":\"counter\",\"delta\":1}\n"

	require.NoError(t, os.WriteFile(path, []byte(history), 0600))

	store := repository.NewStore()

	saver := newTestMetricSaver(t, path, 0, store)
	require.NoError(t, saver.PopulateStore())

	require.Equal(t, int64(6), store.Counters["requests"])
	require.Equal(t, 12.0, store.Gauges["temperature"])

	metricService := NewMetricService(MetricServiceParams{
		Store:  store,
		Saver:  saver,
		Logger: zap.NewNop(),
	})

	metricService.UpdateCounter("requests", 6)
	metricService.UpdateGauge("temperature", 14)

	require.Equal(t, int64(12), store.Counters["requests"])
	require.Equal(t, 14.0, store.Gauges["temperature"])

}

func newTestMetricSaver(
	t *testing.T,
	path string,
	storeInterval int,
	store *repository.MemStorage,
) *MetricSaver {
	t.Helper()

	saver, err := NewMetricSaver(MetricSaverParams{
		Path:          path,
		StoreInterval: storeInterval,
		Store:         store,
		Logger:        zap.NewNop(),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = saver.file.Close()
	})

	return saver
}

func readMetricEvents(t *testing.T, path string) []models.Metrics {
	t.Helper()

	file, err := os.Open(path)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, file.Close())
	}()

	decoder := json.NewDecoder(file)
	events := make([]models.Metrics, 0)

	for i := 0; ; i++ {
		var metric models.Metrics

		err := decoder.Decode(&metric)
		if errors.Is(err, io.EOF) {
			break
		}

		require.NoErrorf(t, err, "Metric decode error %d", i)
		events = append(events, metric)
	}

	return events
}

func requireMetricCounter(
	t *testing.T,
	metric models.Metrics,
	id string,
	delta int64,
) {
	t.Helper()

	require.Equal(t, id, metric.ID)
	require.Equal(t, "counter", metric.MType)
	require.NotNil(t, metric.Delta)
	require.Equal(t, delta, *metric.Delta)
	require.Nil(t, metric.Value)
}

func requireMetricGauge(
	t *testing.T,
	metric models.Metrics,
	id string,
	value float64,
) {
	t.Helper()

	require.Equal(t, id, metric.ID)
	require.Equal(t, "gauge", metric.MType)
	require.NotNil(t, metric.Value)
	require.Equal(t, value, *metric.Value)
	require.Nil(t, metric.Delta)
}
