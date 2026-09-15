-- +goose Up
CREATE SCHEMA IF NOT EXISTS metrics;
CREATE TYPE metrics.metric_type AS ENUM ('counter', 'gauge');
CREATE TABLE IF NOT EXISTS metrics.metrics (
    id    TEXT NOT NULL,
    type  metrics.metric_type NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION,

    CONSTRAINT metrics_pkey
        PRIMARY KEY (id, type)

);
-- +goose Down
DROP TABLE metrics.metrics;

DROP TYPE metrics.metric_type;

DROP SCHEMA IF EXISTS metrics;