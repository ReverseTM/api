package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" envDefault:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address string        `yaml:"address" envDefault:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" envDefault:"4s"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &config
}

func fetchConfigPath() (result string) {
	flag.StringVar(&result, "config", "", "path to config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return
}
