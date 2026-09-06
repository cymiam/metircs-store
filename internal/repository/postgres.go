package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostrgreStorage struct {
	pool *pgxpool.Pool
}

type PostgreStorageParams struct {
	Pool *pgxpool.Pool
}

func NewPostgresStorage(params PostgreStorageParams) *PostrgreStorage {
	return &PostrgreStorage{
		pool: params.Pool,
	}
}

func (p *PostrgreStorage) GetAll(ctx context.Context) ([]models.Metric, error) {
	query, args, err := sq.Select("*").From("metrics.metrics").ToSql()

	if err != nil {
		return nil, fmt.Errorf("cannot create sql query: %w", err)
	}

	rows, err := p.pool.Query(ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("cannot run sql query: %w", err)
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
			return nil, fmt.Errorf("cannot scan metri: %w", err)
		}

		metrics = append(metrics, m)
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return metrics, nil
}

func (p *PostrgreStorage) GetMetric(ctx context.Context, name string, metricType string) (models.Metric, error) {
	query, args, err := sq.Select("id", "type", "delta", "value").
		From("metrics.metrics").
		Where(sq.Eq{"id": name, "type": metricType}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return models.Metric{}, fmt.Errorf("cannot create sql query: %w", err)
	}

	metric := models.Metric{}

	err = p.pool.QueryRow(ctx, query, args...).Scan(
		&metric.ID,
		&metric.MType,
		&metric.Delta,
		&metric.Value,
	)

	if err != nil {
		return models.Metric{}, fmt.Errorf("cannot scan metric: %w", err)
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
		return fmt.Errorf("unknown metric type: %s", metric.MType)
	}
	sql, args, err := builder.ToSql()

	if err != nil {
		return fmt.Errorf("build upsert metric query: %w", err)
	}

	if _, err := p.pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("upsert %s: %w, sql: %s", metric, err, sql)
	}

	return nil

}
