package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexcabrera/justvibin/internal/registry"
)

func TestOpenCmdNotJustvibin(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"open"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(stderr.String(), "Not a justvibin project") {
		t.Fatalf("expected not a project error")
	}
}

func TestOpenCmdProjectNotFound(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"open", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(stderr.String(), "not found") {
		t.Fatalf("expected not found error")
	}
}

func TestOpenCmdSuccess(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, ".justvibin"), []byte(`{"name":"myapp","template":"hypertext","port":59999}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	projectsPath := filepath.Join(baseDir, "justvibin", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"myapp": {Port: 59999, Path: projectDir, Template: "hypertext"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"open"})
	_ = rootCmd.Execute()
}
