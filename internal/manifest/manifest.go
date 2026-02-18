package manifest

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Template struct {
	Name        string
	Description string
	Version     string
	Author      string
	URL         string
}

type Scaffold struct {
	Exclude          []string
	Setup            string
	SetupInteractive bool
}

type Serve struct {
	Type        string
	Dev         string
	Prod        string
	PortEnv     string
	DefaultPort int
	Static      ServeStatic
}

type ServeStatic struct {
	Root       string
	Extensions []string
}

type Project struct {
	MarkerFields []string
}

type Manifest struct {
	Template Template
	Scaffold Scaffold
	Serve    Serve
	Project  Project
}

var namePattern = regexp.MustCompile(`^[a-z0-9-]+$`)

func Parse(data []byte) (Manifest, error) {
	manifest := Manifest{}
	section := ""
	lines := strings.Split(string(data), "\n")
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			if !strings.HasSuffix(line, "]") {
				return Manifest{}, fmt.Errorf("invalid section header at line %d", i+1)
			}
			section = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
			if section == "" {
				return Manifest{}, fmt.Errorf("invalid section header at line %d", i+1)
			}
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return Manifest{}, fmt.Errorf("invalid assignment at line %d", i+1)
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return Manifest{}, fmt.Errorf("invalid assignment at line %d", i+1)
		}
		value := strings.TrimSpace(parts[1])
		if value == "" {
			return Manifest{}, fmt.Errorf("invalid assignment at line %d", i+1)
		}

		switch section {
		case "template":
			assignTemplate(&manifest.Template, key, value)
		case "scaffold":
			assignScaffold(&manifest.Scaffold, key, value)
		case "serve":
			assignServe(&manifest.Serve, key, value)
		case "serve.static":
			assignServeStatic(&manifest.Serve.Static, key, value)
		case "project":
			assignProject(&manifest.Project, key, value)
		}
	}
	return manifest, nil
}

func Validate(manifest Manifest) error {
	var errs []string
	if manifest.Template.Name == "" {
		errs = append(errs, "template.name is required")
	} else if !namePattern.MatchString(manifest.Template.Name) {
		errs = append(errs, "template.name must match [a-z0-9-]+")
	}
	if manifest.Template.Description == "" {
		errs = append(errs, "template.description is required")
	}
	if manifest.Serve.Type == "" {
		errs = append(errs, "serve.type is required")
	} else if manifest.Serve.Type != "static" && manifest.Serve.Type != "command" {
		errs = append(errs, "serve.type must be static or command")
	}
	if manifest.Serve.Type == "command" {
		if manifest.Serve.Dev == "" && manifest.Serve.Prod == "" {
			errs = append(errs, "serve.dev or serve.prod is required for command templates")
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func TemplateName(manifest Manifest) string {
	return manifest.Template.Name
}

func ServeType(manifest Manifest) string {
	return manifest.Serve.Type
}

func ServeCommand(manifest Manifest, mode string) string {
	if mode == "prod" {
		if manifest.Serve.Prod != "" {
			return manifest.Serve.Prod
		}
		return manifest.Serve.Dev
	}
	if manifest.Serve.Dev != "" {
		return manifest.Serve.Dev
	}
	return manifest.Serve.Prod
}

func assignTemplate(template *Template, key, value string) {
	switch key {
	case "name":
		template.Name = trimString(value)
	case "description":
		template.Description = trimString(value)
	case "version":
		template.Version = trimString(value)
	case "author":
		template.Author = trimString(value)
	case "url":
		template.URL = trimString(value)
	}
}

func assignScaffold(scaffold *Scaffold, key, value string) {
	switch key {
	case "exclude":
		scaffold.Exclude = trimArray(value)
	case "setup":
		scaffold.Setup = trimString(value)
	case "setup_interactive":
		scaffold.SetupInteractive = trimBool(value)
	}
}

func assignServe(serve *Serve, key, value string) {
	switch key {
	case "type":
		serve.Type = trimString(value)
	case "dev":
		serve.Dev = trimString(value)
	case "prod":
		serve.Prod = trimString(value)
	case "port_env":
		serve.PortEnv = trimString(value)
	case "default_port":
		serve.DefaultPort = trimInt(value)
	}
}

func assignServeStatic(static *ServeStatic, key, value string) {
	switch key {
	case "root":
		static.Root = trimString(value)
	case "extensions":
		static.Extensions = trimArray(value)
	}
}

func assignProject(project *Project, key, value string) {
	switch key {
	case "marker_fields":
		project.MarkerFields = trimArray(value)
	}
}

func trimString(value string) string {
	value = strings.TrimSpace(value)
	return strings.Trim(value, "\"")
}

func trimBool(value string) bool {
	value = trimString(value)
	return value == "true"
}

func trimInt(value string) int {
	value = trimString(value)
	if value == "" {
		return 0
	}
	var parsed int
	_, err := fmt.Sscanf(value, "%d", &parsed)
	if err != nil {
		return 0
	}
	return parsed
}

func trimArray(value string) []string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := trimString(part)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
