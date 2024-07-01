package config

import (
	"Ecost/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Server  struct {
		Type        string `yaml:"type"`
		BindIP      string `yaml:"bind_ip"`
		Port        string `yaml:"port"`
		MetricsPort string `yaml:"metrics_port"`
		Domain      string `yaml:"domain"`
	} `yaml:"service"`
	Postgresql PostgresConfig `yaml:"postgresql"`
	Redis      RedisConfig    `yaml:"redis"`
	Email      EmailConfig    `yaml:"email"`
	Yandex     YandexConfig   `yaml:"yandex"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type EmailConfig struct {
	Sender   string `json:"sender"`
	Address  string `json:"address"`
	Password string `json:"password"`
}

type YandexConfig struct {
	Bucket string `json:"bucket"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application's configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
