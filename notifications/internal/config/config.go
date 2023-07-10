package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type ConfigKafka struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	Group   string   `yaml:"group"`
}

type ConfigPostgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Config struct {
	Port struct {
		GRPC int `yaml:"grpc"`
		HTTP int `yaml:"http"`
	} `yaml:"port"`

	LogLevel string `yaml:"loglevel"`

	Kafka ConfigKafka `yaml:"kafka"`

	Telegram struct {
		Token  string `yaml:"token"`
		ChatId int64  `yaml:"chat_id"`
	} `yaml:"telegram"`

	Postgres ConfigPostgres `yaml:"postgres"`
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

func (p *ConfigPostgres) URL() string {
	builder := strings.Builder{}
	builder.WriteString("postgres://")
	builder.WriteString(p.User)
	builder.WriteRune(':')
	builder.WriteString(p.Password)
	builder.WriteRune('@')
	builder.WriteString(p.Host)
	builder.WriteRune(':')
	builder.WriteString(strconv.Itoa(p.Port))
	builder.WriteRune('/')
	builder.WriteString(p.Database)
	builder.WriteString("?sslmode=disable")
	return builder.String()
}
