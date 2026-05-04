package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HttpServer  HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Addr        string `yaml:"address" env:"HTTP_SERVER_ADDR" env-required:"true"`
	Timeout     string `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-required:"true"`
	IdleTimeout string `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-required:"true"`
}

func MustLoad() *Config {
	var configPath string

	// flag
	flags := flag.String("config", "", "path to config file")
	flag.Parse()

	// priority: flag > env > default
	if *flags != "" {
		configPath = *flags
	} else if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	} else {
		configPath = defaultConfigPath()
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found : %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file : %s", err.Error)
	}

	return &cfg
}

func defaultConfigPath() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "config/local.yaml"
	}

	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	return filepath.Join(projectRoot, "config", "local.yaml")
}
