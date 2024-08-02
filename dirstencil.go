package main

import (
	"log"
	"os"
	"path"
	"sync"

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
	processTemplateDirForTarget(name, t, templateDir, "")
}

func processTemplateDirForTarget(name string, t target, baseDir string, subDir string) {
	targetDir := path.Join(t.Dir, subDir)
	log.Println("Ensuring dir", targetDir, "exists for", name)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to ensure existence of dir %s for %s: %s", targetDir, name, err)
	}

	dir := path.Join(baseDir, subDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read dir %s for %s: %s", dir, name, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			processTemplateDirForTarget(name, t, baseDir, path.Join(subDir, entry.Name()))
		} else {
			log.Println("Processing", entry.Name(), "in", dir, "for", name)
		}

	}
}

func main() {
	configFile := "dirstencil.yaml"
	log.Println("Loading config from", configFile)
	conf := loadConfig(configFile)
	log.Println("Found configuration with version", conf.Version, "and template dir", conf.TemplateDir)

	wg := sync.WaitGroup{}
	for name, t := range conf.Targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processTarget(name, t, conf.TemplateDir)
		}()
	}
	wg.Wait()

}
