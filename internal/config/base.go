package config

type AppConfig struct {
	APP_VERSION      string
	APP_ENV          string
	JWT_SECRETKEY    string
	JWT_DURATIONTIME string
}
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
	Address string
	Port    string
}

type Config struct {
	APP_CONFIG AppConfig
	Postgres   PostgresDB
	RedisDB    RedisDB
}

type BaseConfig struct {
	APP_VERSION      string `mapstructure:"APP_VERSION"`
	APP_ENV          string `mapstructure:"APP_ENV"`
	JWT_SECRETKEY    string `mapstructure:"JWT_SECRETKEY"`
	JWT_DURATIONTIME string `mapstructure:"JWT_DURATIONTIME"`

	DBuser   string `mapstructure:"POSTGRES_USER"`
	DBName   string `mapstructure:"POSTGRES_DBNAME"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBPort   string `mapstructure:"POSTGRES_PORT"`
	Host     string `mapstructure:"POSTGRES_HOST"`
	DB_URL   string
	SSLMODE  string `mapstructure:"POSTGRES_SSLMODE"`

	RED_ADDR string `mapstructure:"RED_ADDR"`
	RED_PORT string `mapstructure:"RED_PORT"`
}

func (b *BaseConfig) SetupConfig() *Config {
	return &Config{
		APP_CONFIG: AppConfig{
			APP_VERSION:      b.APP_VERSION,
			APP_ENV:          b.APP_ENV,
			JWT_SECRETKEY:    b.JWT_SECRETKEY,
			JWT_DURATIONTIME: b.JWT_DURATIONTIME,
		},
		Postgres: PostgresDB{
			DBuser:   b.DBuser,
			DBName:   b.DBName,
			Port:     b.DBPort,
			Password: b.Password,
			Host:     b.Host,
			SSLMODE:  b.SSLMODE,
			DB_URL:   b.DB_URL,
		},
		RedisDB: RedisDB{
			Address: b.RED_ADDR,
			Port:    b.RED_PORT,
		},
	}
}
