package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" envDefault:"local"`
	JWT        `yaml:"jwt"`
	Storage    `yaml:"storage"`
	HTTPServer `yaml:"http_server"`
}

type JWT struct {
	Secret   string        `yaml:"secret" envDefault:"secret"`
	TokenTTL time.Duration `yaml:"token_ttl" envDefault:"1h"`
}

type Storage struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	SpaceName   string `yaml:"space_name" env-required:"true"`
}

type HTTPServer struct {
	Address string `yaml:"address" envDefault:":8080"`
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
