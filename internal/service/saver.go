package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/internal/repository"
	"go.uber.org/zap"
)

type MetricSaver struct {
	file          *os.File
	storeInterval time.Duration
	store         *repository.MemStorage
	logger        *zap.Logger
	queue         []models.Metrics
}

type MetricSaverParams struct {
	Path          string
	StoreInterval int
	Store         *repository.MemStorage
	Logger        *zap.Logger
}

func NewMetricSaver(params MetricSaverParams) (*MetricSaver, error) {
	file, err := os.OpenFile(params.Path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	return &MetricSaver{
		file:          file,
		storeInterval: time.Duration(params.StoreInterval) * time.Second,
		store:         params.Store,
		logger:        params.Logger,
		queue:         make([]models.Metrics, 0),
	}, nil
}

func (m *MetricSaver) StartTicker() {
	if m.storeInterval <= 0 {
		return
	}

	ticker := time.NewTicker(m.storeInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.WriteToFile(); err != nil {
			m.logger.Error(
				"Cannot marshal store",
				zap.Error(err),
			)

		}
	}
}

func (m *MetricSaver) WriteToFile() error {
	if len(m.queue) == 0 {
		return nil
	}

	encoder := json.NewEncoder(m.file)
	eventsCount := len(m.queue)

	for i := range m.queue {
		if err := encoder.Encode(&m.queue[i]); err != nil {
			// Убираем успешно записанную часть.
			// Текущее и последующие события остаются для повтора.
			m.queue = m.queue[i:]

			return fmt.Errorf(
				"encode pending metric %d: %w",
				i,
				err,
			)
		}
	}

	// Отчистка записанной истории метрик
	m.queue = m.queue[:0]

	if err := m.file.Sync(); err != nil {
		return fmt.Errorf("sync metrics file: %w", err)
	}

	m.logger.Info(
		"Metrics appended to file",
		zap.String("file", m.file.Name()),
		zap.Int("events", eventsCount),
	)

	return nil
}

func (m *MetricSaver) PopulateStore() error {
	if _, err := m.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek metrics file: %w", err)
	}

	decoder := json.NewDecoder(m.file)
	newStore := repository.NewStore()

	metricNumber := 0

	for {
		var metric models.Metrics

		err := decoder.Decode(&metric)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf(
				"Decode metric %d: %w",
				metricNumber,
				err,
			)
		}

		if metric.ID == "" {
			return fmt.Errorf(
				"Metric %d: empty metric ID",
				metricNumber,
			)
		}

		switch metric.MType {
		case "counter":
			if metric.Delta == nil {
				return fmt.Errorf(
					"Metric %d: counter %q has no delta",
					metricNumber,
					metric.ID,
				)
			}

			total := *metric.Delta

			values, ok := newStore.GetCounter(metric.ID)
			if ok && len(values) > 0 {
				total += values[len(values)-1]
			}

			newStore.SetCounter(metric.ID, total)

		case "gauge":
			if metric.Value == nil {
				return fmt.Errorf(
					"Metric %d: gauge %q has no value",
					metricNumber,
					metric.ID,
				)
			}

			newStore.SetGauge(metric.ID, *metric.Value)

		}

		metricNumber++
	}

	// Меняем текущее хранилище только после успешного чтения
	// всего JSONL-файла.
	m.store.Gauges = newStore.Gauges
	m.store.Counters = newStore.Counters

	// // Последующие записи всё равно используют O_APPEND,
	// // но позицию файла также явно перемещаем в конец.
	// if _, err := m.file.Seek(0, io.SeekEnd); err != nil {
	// 	return fmt.Errorf("seek to end of metrics file: %w", err)
	// }

	m.logger.Info(
		"Metrics restored from file",
		zap.String("file", m.file.Name()),
		zap.Int("metric count", metricNumber),
	)

	return nil
}

func (m *MetricSaver) OnMetricChanged(metric models.Metrics) {
	m.queue = append(m.queue, metric)

	// При периодической записи событие заберёт тикер.
	if m.storeInterval > 0 {
		return
	}

	// При interval == 0 запись выполняется синхронно.
	if err := m.WriteToFile(); err != nil {
		if m.logger != nil {
			m.logger.Error(
				"Cannot write metric",
				zap.String("metric", metric.ID),
				zap.Error(err),
			)
		}
	}
}
