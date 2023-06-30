package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigKafka struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	Group   string   `yaml:"group"`
}

type Config struct {
	Port struct {
		GRPC int `yaml:"grpc"`
		HTTP int `yaml:"http"`
	} `yaml:"port"`

	Kafka ConfigKafka `yaml:"kafka"`

	Telegram struct {
		Token  string `yaml:"token"`
		ChatId int64  `yaml:"chat_id"`
	} `yaml:"telegram"`
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
