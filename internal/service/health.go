package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthService struct {
	pool *pgxpool.Pool
}

func NewHealthService(pool *pgxpool.Pool) *HealthService {
	return &HealthService{
		pool: pool,
	}
}

func (service *HealthService) PingDB() error {
	if service.pool != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return service.pool.Ping(ctx)
	}
	return fmt.Errorf("database does not exists")
}
