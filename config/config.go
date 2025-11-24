package config

import (
	"strings"

	"github.com/spf13/viper"
)

type App struct {
	AppPort       string `json:"app_port"`
	AppEnv        string `json:"app_env"`
	PrefixURL     string `json:"prefix_url"`
	ServerTimeOut int    `json:"server_timeout"`
	JwtSecretKey  string `json:"jwt_secret_key"`
	JwtIssuer     string `json:"jwt_issuer"`
	UrlFrontFE    string `json:"url_front_fe"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"db_name"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type Redis struct {
	Addr     string `json:"addr"`
	DB       int    `json:"db"`
	Password string `json:"password"`
}

type Kafka struct {
	Brokers     []string `json:"brokers"`
	TimeoutInMS int      `json:"timeoutInMS"`
	MaxRetry    int      `json:"maxRetry"`
	Topic       string   `json:"topic"`
}

type Minio struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
	UseSSL    bool   `json:"useSSL"`
}

type Config struct {
	App   App    `json:"app"`
	Psql  PsqlDB `json:"psql"`
	Redis Redis  `json:"redis"`
	Kafka Kafka  `json:"kafka"`
	Minio Minio  `json:"minio"`
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort:       viper.GetString("APP_PORT"),
			AppEnv:        viper.GetString("APP_ENV"),
			PrefixURL:     viper.GetString("PREFIX_URL"),
			ServerTimeOut: viper.GetInt("SERVER_TIMEOUT"),
			JwtSecretKey:  viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:     viper.GetString("JWT_ISSUER"),
			UrlFrontFE:    viper.GetString("URL_FRONT_FE"),
		},
		Psql: PsqlDB{
			Host:      viper.GetString("DATABASE_HOST"),
			Port:      viper.GetString("DATABASE_PORT"),
			User:      viper.GetString("DATABASE_USER"),
			Password:  viper.GetString("DATABASE_PASSWORD"),
			DBName:    viper.GetString("DATABASE_NAME"),
			DBMaxOpen: viper.GetInt("DATABASE_MAX_OPEN_CONNECTION"),
			DBMaxIdle: viper.GetInt("DATABASE_MAX_IDLE_CONNECTION"),
		},
		Redis: Redis{
			Addr:     viper.GetString("REDIS_ADDR"),
			DB:       viper.GetInt("REDIS_DB"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
		Kafka: Kafka{
			Brokers:     strings.Split(viper.GetString("KAFKA_BROKERS"), ","),
			TimeoutInMS: viper.GetInt("KAFKA_TIMEOUT_IN_MS"),
			MaxRetry:    viper.GetInt("KAFKA_MAX_RETRY"),
			Topic:       viper.GetString("KAFKA_TOPIC"),
		},
		Minio: Minio{
			Endpoint:  viper.GetString("MINIO_ENDPOINT"),
			AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey: viper.GetString("MINIO_SECRET_KEY"),
			Bucket:    viper.GetString("MINIO_BUCKET"),
			UseSSL:    viper.GetBool("MINIO_USE_SSL"),
		},
	}
}
