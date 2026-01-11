package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port     string `yaml:"port"`
	LogLevel string `yaml:"log_level"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type Config struct {
	App      AppConfig      `yaml:"app"`
	Postgres PostgresConfig `yaml:"postgres"`
	Telegram TelegramConfig `yaml:"telegram"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
