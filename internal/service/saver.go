package service

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/cymiam/metircs-store/internal/logger"
	"github.com/cymiam/metircs-store/internal/repository"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type MetricSaver struct {
	file          *os.File
	storeInterval time.Duration
	store         *repository.MemStorage
	reader        *bufio.Reader
	writer        *bufio.Writer
}

type MetricSaverParams struct {
	Path          string
	StoreInterval int
	Store         *repository.MemStorage
}

func NewMetricSaver(params MetricSaverParams) MetricSaver {
	file, err := os.OpenFile(params.Path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		logger.Log.Fatal("Cannot open file", zap.String("name: ", params.Path))
	}

	return MetricSaver{
		file:          file,
		storeInterval: time.Duration(params.StoreInterval) * time.Second,
		store:         params.Store,
		reader:        bufio.NewReader(file),
		writer:        bufio.NewWriter(file),
	}
}

func (m *MetricSaver) Update() {
	if m.storeInterval <= 0 {
		return
	}

	ticker := time.NewTicker(m.storeInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.WriteToFile(); err != nil {
			logger.Log.Error(
				"Cannot marshal store",
				zap.Error(err),
			)
		}
	}
}

func (m *MetricSaver) WriteToFile() error {
	if err := m.writer.Flush(); err != nil {
		return err
	}

	if err := m.file.Truncate(0); err != nil {
		return err
	}

	if _, err := m.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	m.writer.Reset(m.file)

	written, err := easyjson.MarshalToWriter(m.store, m.writer)
	if err != nil {
		return err
	}
	logger.Log.Info("Written into file", zap.String("file name", m.file.Name()), zap.Int("Bytes", written))

	m.writer.Flush()
	return nil
}

func (m *MetricSaver) PopulateStore() error {

	if _, err := m.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	m.reader.Reset(m.file)

	err := easyjson.UnmarshalFromReader(m.reader, m.store)
	if err != nil {
		return err
	}

	logger.Log.Info("Read from file", zap.String("file name", m.file.Name()))
	return nil
}

func (m *MetricSaver) OnMetricChanged() {
	if m.storeInterval != 0 {
		return
	}

	if err := m.WriteToFile(); err != nil {
		logger.Log.Error("Cannot write to file", zap.Error(err))
	}
}
