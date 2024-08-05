package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"text/template"

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

type renderer struct {
	config configuration
}

func (r *renderer) loadConfig(configFile string) error {
	log.Println("Loading config from", configFile)
	file, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", configFile, err)
	}
	var conf configuration
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return fmt.Errorf("failed to unmarshall file %s: %w", configFile, err)
	}
	log.Println("Found configuration with version", conf.Version, "and template dir", conf.TemplateDir)
	r.config = conf
	return nil
}

func (r *renderer) processTarget(name string, t target) error {
	log.Println("Processing target", name)
	if err := r.processTemplateDirForTarget(name, t, r.config.TemplateDir, ""); err != nil {
		return fmt.Errorf("failed to process target %s: %w", name, err)
	}
	return nil
}

func (r *renderer) processTemplateDirForTarget(name string, t target, baseDir string, subDir string) error {
	targetDir := path.Join(t.Dir, subDir)
	log.Println("Ensuring dir", targetDir, "exists for", name)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to ensure existence of dir %s: %w", targetDir, err)
	}

	dir := path.Join(baseDir, subDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read dir %s: %w", dir, err)
	}

	for _, entry := range entries {
		log.Println("Processing", entry.Name(), "in", dir, "for", name)
		if entry.IsDir() {
			if err := r.processTemplateDirForTarget(name, t, baseDir, path.Join(subDir, entry.Name())); err != nil {
				// don't wrap error for recursive call
				return err
			}
		} else {
			if err := r.renderTemplateForTarget(name, t, baseDir, subDir, entry.Name()); err != nil {
				return fmt.Errorf("failed to render %s: %w", entry.Name(), err)
			}
		}
	}
	return nil
}

func (r *renderer) renderTemplateForTarget(name string, t target, baseDir string, subDir string, file string) error {
	targetFile := path.Join(t.Dir, subDir, file)
	sourceFile := path.Join(baseDir, subDir, file)
	log.Println("Rendering", sourceFile, "to", targetFile, "for", name)

	tpl, err := template.New(file).
		Delims(r.config.Delims.Left, r.config.Delims.Right).
		ParseFiles(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", sourceFile, err)
	}

	values := t.Vars
	values["genmsg"] = "Generated by treestencil"
	buffer := bytes.Buffer{}
	if err := tpl.Execute(&buffer, values); err != nil {
		return fmt.Errorf("failed to render template %s: %s", sourceFile, err)
	}

	f, err2 := os.Create(targetFile)
	if err2 != nil {
		return fmt.Errorf("failed create target file %s: %s", targetFile, err)
	}
	defer f.Close()
	if _, err := f.WriteString(buffer.String()); err != nil {
		return fmt.Errorf("failed to write rendered template to %s: %w", targetFile, err)
	}
	return nil
}

func (r *renderer) renderAll() error {
	log.Println("Rendering for", len(r.config.Targets), "targets")
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(r.config.Targets))
	for name, t := range r.config.Targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- r.processTarget(name, t)
		}()
	}
	wg.Wait()
	var errsFound = false
	for range r.config.Targets {
		err := <-errChan
		if err != nil {
			errsFound = true
			log.Printf("Failure: %s", err)
		}
	}
	if errsFound {
		return fmt.Errorf("encountered errors")
	}
	return nil
}

func NewRenderer(configFile string) (*renderer, error) {
	r := renderer{}
	if err := r.loadConfig(configFile); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return &r, nil
}

func main() {
	configFile := "treestencil.yaml"
	r, err := NewRenderer(configFile)
	if err != nil {
		log.Fatalf("Initialization failed: %s", err)
	}
	if err := r.renderAll(); err != nil {
		log.Fatalf("Failed: %s", err)
	}
}