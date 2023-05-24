package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigService struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Config struct {
	Port int `yaml:"port"`

	Services struct {
		Loms ConfigService `yaml:"loms"`
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
