package manifest

import (
	"strings"
	"testing"
)

func TestParseAndValidateManifest(t *testing.T) {
	input := strings.Join([]string{
		"[template]",
		"name = \"hypertext\"",
		"description = \"HTML + HTMX\"",
		"version = \"1.0.0\"",
		"",
		"[scaffold]",
		"exclude = [\".git\", \".DS_Store\"]",
		"setup = \"./setup.sh\"",
		"setup_interactive = true",
		"",
		"[serve]",
		"type = \"command\"",
		"dev = \"./start.sh --dev\"",
		"prod = \"./start.sh\"",
		"port_env = \"PORT\"",
		"default_port = 8000",
		"",
		"[serve.static]",
		"root = \".\"",
		"extensions = [\".html\"]",
		"",
		"[project]",
		"marker_fields = [\"template\", \"port\"]",
	}, "\n")

	manifest, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if manifest.Template.Name != "hypertext" || manifest.Template.Description != "HTML + HTMX" {
		t.Fatalf("unexpected template values")
	}
	if manifest.Scaffold.Setup != "./setup.sh" || !manifest.Scaffold.SetupInteractive {
		t.Fatalf("unexpected scaffold values")
	}
	if manifest.Serve.Type != "command" || manifest.Serve.Dev == "" || manifest.Serve.Prod == "" {
		t.Fatalf("unexpected serve values")
	}
	if manifest.Serve.DefaultPort != 8000 {
		t.Fatalf("unexpected default port")
	}
	if len(manifest.Scaffold.Exclude) != 2 || manifest.Scaffold.Exclude[0] != ".git" {
		t.Fatalf("unexpected exclude list")
	}
	if len(manifest.Serve.Static.Extensions) != 1 || manifest.Serve.Static.Extensions[0] != ".html" {
		t.Fatalf("unexpected static extensions")
	}
	if len(manifest.Project.MarkerFields) != 2 {
		t.Fatalf("unexpected marker fields")
	}

	if err := Validate(manifest); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if TemplateName(manifest) != "hypertext" {
		t.Fatalf("unexpected template name helper")
	}
	if ServeType(manifest) != "command" {
		t.Fatalf("unexpected serve type helper")
	}
	if ServeCommand(manifest, "dev") == "" || ServeCommand(manifest, "prod") == "" {
		t.Fatalf("expected serve command helpers")
	}
}

func TestValidateManifestErrors(t *testing.T) {
	manifest := Manifest{}
	if err := Validate(manifest); err == nil {
		t.Fatalf("expected validation error")
	}

	manifest.Template.Name = "Invalid"
	manifest.Template.Description = "desc"
	manifest.Serve.Type = "bad"
	if err := Validate(manifest); err == nil {
		t.Fatalf("expected validation error")
	}

	manifest.Template.Name = "valid-name"
	manifest.Serve.Type = "command"
	manifest.Serve.Dev = ""
	manifest.Serve.Prod = ""
	if err := Validate(manifest); err == nil {
		t.Fatalf("expected validation error for command serve")
	}
}

func TestParseInvalidManifestReturnsError(t *testing.T) {
	input := "[template\nname = \"oops\"\n"
	if _, err := Parse([]byte(input)); err == nil {
		t.Fatalf("expected parse error")
	}
}
