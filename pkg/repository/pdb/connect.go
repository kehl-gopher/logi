package pdb

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

type Database interface {
	ConnectPostgres() error
	PingDB() bool
	Close()
	DB() *bun.DB
	Insert(ctx context.Context, model interface{}) error
	InsertMany(ctx context.Context, models ...interface{}) error
	CheckExists(ctx context.Context, query string, model interface{}) (bool, error)
	SelectSingle(ctx context.Context, model interface{}, query string, args ...interface{}) error
	UpdateModel(ctx context.Context, model interface{}, column string, query string, args ...interface{}) error
}

type postgresConn struct {
	conf *config.Config
	log  *utils.Log
	bun  *bun.DB
}

func NewPostgresConn(conf *config.Config, logs *utils.Log) Database {
	return &postgresConn{
		conf: conf,
		log:  logs,
	}
}

func (p *postgresConn) ConnectPostgres() error {
	port, err := utils.PortResolver(p.conf.Postgres.Port)

	if err != nil {
		utils.PrintLog(p.log, err.Error(), utils.ErrorLevel)
		return err
	}

	dsn := connString(
		p.conf.Postgres.DBuser,
		p.conf.Postgres.Password,
		p.conf.Postgres.Host,
		p.conf.Postgres.DBName, p.conf.Postgres.SSLMODE, port)

	connPool, err := poolConfig(dsn)
	if err != nil {
		utils.PrintLog(p.log, err.Error(), utils.ErrorLevel)
		return err
	}

	sqldb := stdlib.OpenDBFromPool(connPool)
	bun := bun.NewDB(sqldb, pgdialect.New())

	if err := bun.Ping(); err != nil {
		utils.PrintLog(p.log, err.Error(), utils.FatalLevel)
		return err
	}
	p.bun = bun

	utils.PrintLog(p.log, "Database connection successful", utils.InfoLevel)
	return nil
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

func (p *postgresConn) PingDB() bool {
	return p.bun.Ping() == nil
}

func (p *postgresConn) Close() {
	if p.bun != nil {
		p.bun.Close()
	}
}

func (p *postgresConn) DB() *bun.DB {
	return p.bun
}
