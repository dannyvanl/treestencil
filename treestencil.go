package main

import (
	"bytes"
	"fmt"
	"log"
	"maps"
	"os"
	"path"
	"sync"
	"text/template"
)

type treestencil struct {
	config configuration
}

func (ts *treestencil) processTarget(name string, t target) error {
	log.Println("Processing target", name)
	if err := ts.processTemplateDirForTarget(name, t, ts.config.TemplateDir, ""); err != nil {
		return fmt.Errorf("failed to process target %s: %w", name, err)
	}
	return nil
}

func (ts *treestencil) processTemplateDirForTarget(name string, t target, baseDir string, subDir string) error {
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
			if err := ts.processTemplateDirForTarget(name, t, baseDir, path.Join(subDir, entry.Name())); err != nil {
				// don't wrap error for recursive call
				return err
			}
		} else {
			if err := ts.renderTemplateForTarget(name, t, baseDir, subDir, entry.Name()); err != nil {
				return fmt.Errorf("failed to render %s: %w", entry.Name(), err)
			}
		}
	}
	return nil
}

func (ts *treestencil) renderTemplateForTarget(name string, t target, baseDir string, subDir string, file string) error {
	targetFile := path.Join(t.Dir, subDir, file)
	sourceFile := path.Join(baseDir, subDir, file)
	log.Println("Rendering", sourceFile, "to", targetFile, "for", name)

	tpl, err := template.New(file).
		Delims(ts.config.Delims.Left, ts.config.Delims.Right).
		ParseFiles(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", sourceFile, err)
	}

	values := make(map[string]interface{})
	maps.Copy(values, ts.config.Vars)
	maps.Copy(values, t.Vars)

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

func (ts *treestencil) stencilAll() error {
	log.Println("Stenciling for", len(ts.config.Targets), "targets")
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(ts.config.Targets))
	for name, t := range ts.config.Targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- ts.processTarget(name, t)
		}()
	}
	wg.Wait()
	var errsFound = false
	for range ts.config.Targets {
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

func newTreestencil(config configuration) *treestencil {
	return &treestencil{config: config}
}

func main() {
	configFile := "treestencil.yaml"
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	ts := newTreestencil(*config)
	if err := ts.stencilAll(); err != nil {
		log.Fatalf("Failed to run: %s", err)
	}
}
