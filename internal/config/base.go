package config

type RabbitMQ struct {
	CONN_STR string
}
type AppConfig struct {
	APP_VERSION      string
	APP_ENV          string
	APP_URL          string
	FRONTEND_URL     string
	JWT_SECRETKEY    string
	JWT_DURATIONTIME string
	SMTP_USERNAME    string
	SMTP_PORT        string
	SMTP_HOST        string
	SMTP_PASSWORD    string
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
	RabbitMQ   RabbitMQ
}

type BaseConfig struct {
	// application config
	APP_VERSION  string `mapstructure:"APP_VERSION"`
	APP_ENV      string `mapstructure:"APP_ENV"`
	APP_URL      string `mapstructure:"APP_URL"`
	FRONTEND_URL string `mapstructure:"FRONTEND_URL"`

	// jwt secret key
	JWT_SECRETKEY    string `mapstructure:"JWT_SECRETKEY"`
	JWT_DURATIONTIME string `mapstructure:"JWT_DURATIONTIME"`

	// smtp config
	SMTP_USERNAME string `mapstructure:"SMTP_USERNAME"`
	SMTP_PORT     string `mapstructure:"SMTP_PORT"`
	SMTP_HOST     string `mapstructure:"SMTP_HOST"`
	SMTP_PASSWORD string `mapstrucutre:"SMTP_PASSWORD"`

	// database  config
	DBuser   string `mapstructure:"POSTGRES_USER"`
	DBName   string `mapstructure:"POSTGRES_DBNAME"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBPort   string `mapstructure:"POSTGRES_PORT"`
	Host     string `mapstructure:"POSTGRES_HOST"`
	DB_URL   string
	SSLMODE  string `mapstructure:"POSTGRES_SSLMODE"`

	// redis config
	RED_ADDR string `mapstructure:"RED_ADDR"`
	RED_PORT string `mapstructure:"RED_PORT"`

	// rabbitmq config
	Rabbit_MQ string `mapstructure:"RABBIT_MQ"`
}

func (b *BaseConfig) SetupConfig() *Config {
	return &Config{
		APP_CONFIG: AppConfig{
			APP_VERSION:      b.APP_VERSION,
			APP_ENV:          b.APP_ENV,
			JWT_SECRETKEY:    b.JWT_SECRETKEY,
			JWT_DURATIONTIME: b.JWT_DURATIONTIME,
			SMTP_USERNAME:    b.SMTP_USERNAME,
			SMTP_PORT:        b.SMTP_PORT,
			SMTP_PASSWORD:    b.SMTP_PASSWORD,
			SMTP_HOST:        b.SMTP_HOST,
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
		RabbitMQ: RabbitMQ{
			CONN_STR: b.Rabbit_MQ,
		},
	}
}
