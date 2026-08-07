package models

type Metric[T float64 | int64] struct {
	Value T
	Name  string
}

func NewMetric[T float64 | int64](name string, value T) Metric[T] {
	return Metric[T]{
		Value: value,
		Name:  name,
	}
}
