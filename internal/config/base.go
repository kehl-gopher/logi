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
	Address string
	Port    string
}

type Config struct {
	Postgres PostgresDB
	RedisDB  RedisDB
}

type BaseConfig struct {
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
