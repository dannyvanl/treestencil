package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Version int64 `yaml:"version"`
}

func loadConfig(configFile string) configuration {
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Failed to read configuration from ", configFile, err)
	}
	var conf configuration
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("Failed to unmarshall configuraton from ", configFile, err)
	}
	return conf
}

func main() {
	configFile := "dirstencil.yaml"
	log.Println("Loading config from", configFile)
	conf := loadConfig(configFile)
	log.Println("Found configuration with version", conf.Version)

}
