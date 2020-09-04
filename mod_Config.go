package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const bmExcludePathsFile = "BMexclude.txt"

// Config struct
type Config struct {
	Service struct {
		Port int    `yaml:"port"`
		Url  string `yaml:"url"`
	} `yaml:"service"`
	Paths struct {
		Distrs string `yaml:"distrs"`
	} `yaml:"paths"`
}

func readConfigFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func readConfigEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatal(err)
	}
}