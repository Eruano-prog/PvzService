package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Address string        `yaml:"address" env:"HTTP_ADDRESS" env-default:":8080"`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"5s"`
}

type Config struct {
	LogLevel          string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	HTTPConfig        HTTPConfig    `yaml:"api_server"`
	GrpcAddress       string        `yaml:"grpc_address" env:"GRPC_ADDRESS" env-default:":3000"`
	PrometheusAddress string        `yaml:"prometheus_address" env:"PROMETHEUS_ADDRESS" env-default:":9000"`
	DBAddress         string        `yaml:"db_address" env:"DB_ADDRESS" env-default:"localhost:5432"`
	TokenTTL          time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"24h"`
	TokenSecret       string        `yaml:"token_secret" env:"TOKEN_SECRET" env-default:"testingParam"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return cfg
}
