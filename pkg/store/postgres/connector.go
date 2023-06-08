package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionConfig struct {
	Host               string
	Port               int
	User               string
	Password           string
	DbName             string
	MaxOpenConnections int
	MaxIdleTime        string
}

func OpenDB(cfg ConnectionConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	pgxConfig.MaxConns = int32(cfg.MaxOpenConnections)
	pgxConfig.MaxConnIdleTime = duration

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
