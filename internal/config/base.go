package config

type RabbitMQ struct {
	CONN_STR string
}
type AppConfig struct {
	APP_VERSION  string
	APP_ENV      string
	APP_URL      string
	FRONTEND_URL string

	// jwt tokens
	JWT_SECRETKEY    string
	JWT_DURATIONTIME string

	// smtp conf
	SMTP_USERNAME string
	SMTP_PORT     string
	SMTP_HOST     string
	SMTP_PASSWORD string

	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_API_KEY       string
	GOOGLE_CALLBACK      string
	OAUTH_PASSWORD       string
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

	// google client config
	GOOGLE_CLIENT_ID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GOOGLE_CLIENT_SECRET string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GOOGLE_API_KEY       string `mapstructure:"GOOGLE_API_KEY"`
	GOOGLE_CALLBACK      string `mapstructure:"GOOGLE_CALLBACK"`
	OAUTH_PASSWORD       string `mapstructure:"OAUTH_PASSWORD"`
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
			APP_VERSION: b.APP_VERSION,
			APP_ENV:     b.APP_ENV,
			APP_URL:     b.APP_URL,

			FRONTEND_URL: b.FRONTEND_URL,

			JWT_SECRETKEY:    b.JWT_SECRETKEY,
			JWT_DURATIONTIME: b.JWT_DURATIONTIME,

			SMTP_USERNAME: b.SMTP_USERNAME,
			SMTP_PORT:     b.SMTP_PORT,
			SMTP_PASSWORD: b.SMTP_PASSWORD,
			SMTP_HOST:     b.SMTP_HOST,

			GOOGLE_CLIENT_ID:     b.GOOGLE_CLIENT_ID,
			GOOGLE_API_KEY:       b.GOOGLE_API_KEY,
			GOOGLE_CLIENT_SECRET: b.GOOGLE_CLIENT_SECRET,
			GOOGLE_CALLBACK:      b.GOOGLE_CALLBACK,
			OAUTH_PASSWORD:       b.OAUTH_PASSWORD,
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
