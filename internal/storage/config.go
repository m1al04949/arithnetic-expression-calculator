package storage

type Config struct {
	ConfigPath  string
	DatabaseURL string `yaml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
