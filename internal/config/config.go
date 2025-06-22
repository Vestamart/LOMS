package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ClientConfig struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

type gRPCServerConfig struct {
	Port string `yaml:"gRPCport"`
}

type HTTPServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	LOMSServer gRPCServerConfig `yaml:"loms_server"`
	Database   DatabaseConfig   `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
