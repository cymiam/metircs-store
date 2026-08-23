package logger

import (
	"go.uber.org/zap"
)

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func NewLogger(level, name string) (*zap.Logger, error) {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return log.With(zap.String("Serivce", name)), nil
}
