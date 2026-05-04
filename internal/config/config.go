package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string     `yaml:"storagePath" env:"STORAGE_PATH" env-required:"true"`
	HttpServer  HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Addr        string `yaml:"addr" env:"HTTP_SERVER_ADDR" env-required:"true"`
	Timeout     string `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-required:"true"`
	IdleTimeout string `yaml:"idleTimeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-required:"true"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to config file")

		flag.Parse()

		configPath = *flags

		if configPath == "" {
			panic("config path is required")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found : %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("failed to read config file : %s", err.Error())
	}

	return &cfg
}
