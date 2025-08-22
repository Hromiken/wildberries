package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"log"`
		PG    `yaml:"postgres"`
		Kafka `yaml:"kafka"`
		Cache `yaml:"cache"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		MaxPoolSize int    `env-required:"true" yaml:"max_pool_size" env:"PG_MAX_POOL_SIZE"`
		URL         string `env-required:"true"                      env:"PG_URL"`
	}

	Kafka struct {
		Brokers []string `env:"KAFKA_BROKERS" env-separator:"," env-required:"true"`
		Topic   string   `env:"KAFKA_TOPIC" env-default:"orders"`
		GroupID string   `env:"KAFKA_GROUP_ID" env-default:"order-consumer-group"`
	}

	Cache struct {
		CacheSize int           `env-required:"true" yaml:"size" env:"CACHE_SIZE"`
		TTL       time.Duration `env-required:"true" yaml:"ttl" env:"CACHE_TTL"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	_ = godotenv.Load()

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
