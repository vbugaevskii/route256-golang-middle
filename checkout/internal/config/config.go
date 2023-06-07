package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigService struct {
	Netloc string `yaml:"netloc"`
	Token  string `yaml:"token"`
}

type Config struct {
	Port struct {
		GRPC int `yaml:"grpc"`
		HTTP int `yaml:"http"`
	} `yaml:"port"`

	Services struct {
		Loms           ConfigService `yaml:"loms"`
		ProductService ConfigService `yaml:"product_service"`
	} `yaml:"services"`
}

var AppConfig = Config{}

func Init() error {
	rawYaml, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	err = yaml.Unmarshal(rawYaml, &AppConfig)
	if err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}