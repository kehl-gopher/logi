package config

type PostgresDB struct {
	DBuser   string
	DBName   string
	Password string
	Port     string
	Host     string
	DB_URL   string
	SSLMODE  string
}

type RedisDB struct {
	RedisUrl string
}

type Config struct {
	Postgres PostgresDB
	RedisDB  RedisDB
}
