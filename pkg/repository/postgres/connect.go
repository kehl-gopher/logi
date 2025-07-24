package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type PostgresDB interface {
	ConnectPostgres() (*bun.DB, error)
}

type postgresConn struct {
	conf *config.Config
	log  *utils.Log
}

func NewPostgresConn(conf *config.Config, logs *utils.Log) PostgresDB {
	return &postgresConn{
		conf: conf,
		log:  logs,
	}
}

func (p *postgresConn) ConnectPostgres() (*bun.DB, error) {
	port, err := utils.PortResolver(p.conf.Postgres.Port)

	if err != nil {
		utils.PrintLog(p.log, err.Error(), utils.ErrorLevel)
		return nil, err
	}

	dsn := connString(
		p.conf.Postgres.DBuser,
		p.conf.Postgres.Password,
		p.conf.Postgres.Host,
		p.conf.Postgres.DBName, p.conf.Postgres.SSLMODE, port)

	connPool, err := poolConfig(dsn)
	if err != nil {
		utils.PrintLog(p.log, err.Error(), utils.ErrorLevel)
		return nil, err
	}

	sqldb := stdlib.OpenDBFromPool(connPool)
	bun := bun.NewDB(sqldb, pgdialect.New())
	return bun, nil
}

func poolConfig(connString string) (*pgxpool.Pool, error) {
	const (
		maxConn           int32 = 10
		minConn           int32 = 1
		maxConnIdleTime         = 10 * time.Minute
		maxConnLifetime         = 30 * time.Minute
		healthCheckPeriod       = 10 * time.Second
		connectTimeout          = 5 * time.Second
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	config.MinConns = minConn
	config.MaxConns = maxConn
	config.MaxConnIdleTime = maxConnIdleTime
	config.MaxConnLifetime = maxConnLifetime
	config.HealthCheckPeriod = healthCheckPeriod
	config.ConnConfig.ConnectTimeout = connectTimeout

	p, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func connString(user, password, host, dbname, sslmode string, port int) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}
