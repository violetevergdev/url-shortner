package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Env         string   `yaml:"env" env-default:"local"`
	StoragePath string   `yaml:"storage_path" env-required:"true"`
	HTTPServer  HttpConf `yaml:"http_server"`
}

type HttpConf struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func NewConfig() *Config {
	var c Config

	buf, err := os.ReadFile("config/local.yaml")
	if err != nil {
		log.Fatal("Config file not found")
	}
	_ = yaml.Unmarshal([]byte(buf), &c)

	return &c
}
