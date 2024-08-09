package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type delims struct {
	Left  string `yaml:"left"`
	Right string `yaml:"right"`
}

type target struct {
	Dir  string                 `yaml:"dir"`
	Vars map[string]interface{} `yaml:"vars"`
}

type configuration struct {
	Version     int64             `yaml:"version"`
	TemplateDir string            `yaml:"template-dir"`
	Delims      delims            `yaml:"delims"`
	Targets     map[string]target `yaml:"targets"`
}

func loadConfig(configFile string) (*configuration, error) {
	log.Println("Loading config from", configFile)
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", configFile, err)
	}
	var conf configuration
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall file %s: %w", configFile, err)
	}
	log.Println("Found configuration with version", conf.Version, "and template dir", conf.TemplateDir)
	return &conf, nil
}
