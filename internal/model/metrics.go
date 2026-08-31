package models

import (
	"fmt"
	"strconv"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metric) String() string {
	return "Metric(" + m.MType + "):" + m.ID
}

func (m Metric) MetricValue() (string, error) {

	switch m.MType {

	case "counter":
		if m.Delta == nil {
			return "", fmt.Errorf("%s не содержит delta", m)
		}

		return strconv.FormatInt(*m.Delta, 10), nil

	case "gauge":
		if m.Value == nil {
			return "", fmt.Errorf("%s не содержит value", m)
		}

		return strconv.FormatFloat(
			*m.Value,
			'f',
			-1,
			64,
		), nil
	}

	return "", fmt.Errorf(
		"неизвестный тип метрики %q",
		m.MType,
	)
}
