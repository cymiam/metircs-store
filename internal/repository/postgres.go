package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostrgreStorage struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

type PostgreStorageParams struct {
	Pool   *pgxpool.Pool
	Logger *zap.Logger
}

func NewPostgresStorage(params PostgreStorageParams) *PostrgreStorage {
	return &PostrgreStorage{
		pool:   params.Pool,
		logger: params.Logger,
	}
}

func (p *PostrgreStorage) GetAll(ctx context.Context) ([]models.Metric, error) {
	query, args, err := sq.Select("*").From("metrics.metrics").ToSql()

	if err != nil {
		p.logger.Error("cannot create sql query", zap.Error(err))
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)

	if err != nil {
		p.logger.Error("cannot run sql query", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	metrics := make([]models.Metric, 0)

	for rows.Next() {
		var m models.Metric
		err := rows.Scan(
			&m.ID,
			&m.MType,
			&m.Delta,
			&m.Value,
		)

		if err != nil {
			p.logger.Error("cannot scan metric", zap.Error(err))
			return nil, err
		}

		metrics = append(metrics, m)
	}

	err = rows.Err()

	if err != nil {
		p.logger.Error("error in rows", zap.Error(err))
		return nil, err
	}

	return metrics, nil
}

func (p *PostrgreStorage) GetMetric(ctx context.Context, name string, metricType string) (models.Metric, error) {
	query, args, err := sq.Select("*").From("metrics.metrics").ToSql()

	if err != nil {
		p.logger.Error("cannot create sql query", zap.Error(err))
		return models.Metric{}, err
	}

	rows, err := p.pool.Query(ctx, query, args...)

	if err != nil {
		p.logger.Error("cannot run sql query", zap.Error(err))
		return models.Metric{}, err
	}

	defer rows.Close()

	metric := models.Metric{}

	err = rows.Scan(
		&metric.ID,
		&metric.MType,
		&metric.Delta,
		&metric.Value,
	)

	if err != nil {
		p.logger.Error("cannot scan metric", zap.Error(err))
		return models.Metric{}, err
	}

	err = rows.Err()

	if err != nil {
		p.logger.Error("error in rows", zap.Error(err))
		return models.Metric{}, err
	}

	return metric, nil
}

func (p *PostrgreStorage) SetMetric(ctx context.Context, metric models.Metric) error {

	var builder sq.InsertBuilder

	switch metric.MType {
	case "counter":
		// Upsert metric into table metrics (if metric exists, update its delta)
		builder = sq.Insert("metrics.metrics AS current").
			PlaceholderFormat(sq.Dollar).
			Columns("id", "type", "delta").
			Values(metric.ID, metric.MType, metric.Delta).
			Suffix(`ON CONFLICT (id, type)
				    DO UPDATE SET delta = current.delta + EXCLUDED.delta`)
	case "gauge":
		// Upsert metric into table metrics (if metric exists, update its value)
		builder = sq.Insert("metrics.metrics").
			PlaceholderFormat(sq.Dollar).
			Columns("id", "type", "value").
			Values(metric.ID, metric.MType, metric.Value).
			Suffix(`ON CONFLICT (id, type)
				    DO UPDATE SET value = EXCLUDED.value`)
	default:
		return fmt.Errorf("Unknown metric type: %s", metric.MType)
	}
	sql, args, err := builder.ToSql()

	if err != nil {
		return fmt.Errorf("build upsert metric query: %w", err)
	}

	if _, err := p.pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("upsert %s: err%w, sql: %s", metric, err, sql)
	}

	return nil

}
