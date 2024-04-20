package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	DatabaseURL string        `yaml:"database_url" env-required:"true"`
	Workers     int           `yaml:"workers"`
	Token       string        `yaml:"token" env-required:"true"`
	TokenExpire time.Duration `yaml:"token_expire"`
	HTTPServer  `yaml:"http_server"`
	Timeouts
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Timeouts struct {
	ProcessingInterval   time.Duration `yaml:"process_interval"`
	OperationSumInterval time.Duration `yaml:"operation_sum"`
	OperationSubInterval time.Duration `yaml:"operation_sub"`
	OperationMulInterval time.Duration `yaml:"operation_mul"`
	OperationDivInterval time.Duration `yaml:"operation_div"`
}

func MustLoadCfg() *Config {

	// Path to YAML file in project's directory
	path := "/config/config.yaml"

	// Get Config Path
	currDir, _ := os.Getwd()
	projectDir := filepath.Join(currDir, "..", "..")
	configPath := filepath.Join(projectDir, path)
	if configPath == "" {
		panic("config path is not set")
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	// Reading config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
