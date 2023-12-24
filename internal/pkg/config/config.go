package config

import (
	"log/slog"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port int    `envconfig:"PORT"`
	Host string `envconfig:"HOST"`

	LogLevel slog.Level `envconfig:"LOG_LEVEL"`

	DBConnectionString  string        `envconfig:"DB_CONNECTION_STRING"`
	DBConnectionTimeout time.Duration `envconfig:"DB_CONNECTION_TIMEOUT"`
	DBName              string        `envconfig:"DB_NAME"`

	// S3
	S3BucketName      string `envconfig:"S3_BUCKET_NAME"`
	S3AWSRegion       string `envconfig:"S3_AWS_REGION"`
	S3AccessKeyID     string `envconfig:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey string `envconfig:"S3_SECRET_ACCESS_KEY"`
	S3UseSSL          bool   `envconfig:"S3_USE_SSL"`
}

func new(c Config) (*Config, error) {
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func Must(c Config) *Config {

	conf, err := new(c)

	if err != nil {
		panic(err)
	}

	return conf
}
