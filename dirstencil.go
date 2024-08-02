package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type target struct {
	Dir  string            `yaml:"dir"`
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
		log.Fatalf("Failed to read configuration from %s: %s", configFile, err)
	}
	var conf configuration
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatalf("Failed to unmarshall configuraton from %s: %s", configFile, err)
	}
	return conf
}

func processTarget(name string, t target, templateDir string) {
	log.Println("Processing target", name)

	log.Println("Ensuring dir", t.Dir, "exists for", name)
	if err := os.MkdirAll(t.Dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to ensure existence of dir %s: %s", t.Dir, err)
	}
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
