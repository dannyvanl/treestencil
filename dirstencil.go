package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type target struct {
	Vars map[string]string `yaml:"vars"`
}

type configuration struct {
	Version     int64             `yaml:"version"`
	TemplateDir string            `yaml:"template-dir"`
	Targets     map[string]target `yaml:"targets"`
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

func processTarget(name string, t target, templateDir string) {
	log.Println("Processing target", name)
}

func main() {
	configFile := "dirstencil.yaml"
	log.Println("Loading config from", configFile)
	conf := loadConfig(configFile)
	log.Println("Found configuration with version", conf.Version, "and template dir", conf.TemplateDir)

	for name, target := range conf.Targets {
		processTarget(name, target, conf.TemplateDir)
	}

}
