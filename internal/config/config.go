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

type Config struct {
	LOMSServer gRPCServerConfig `yaml:"loms_server"`
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
