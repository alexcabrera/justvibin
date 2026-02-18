package config

import (
	"os"
	"path/filepath"
)

const (
	ConfigDirName        = "justvibin"
	ProjectsFileName     = "projects.json"
	TemplatesDirName     = "templates"
	TemplatesFileName    = "templates.toml"
	CaddyfileName        = "Caddyfile"
	ConfigFileName       = "config.toml"
	ProxyLogName         = "proxy.log"
	ProxyErrName         = "proxy.err"
	ProxyLabel           = "land.charm.justvibin.proxy"
	BasePort             = 3000
	DefaultTemplatesPath = "~/.config/justvibin/templates.toml"
)

const launchAgentsDir = "Library/LaunchAgents"

func ConfigDir() (string, error) {
	base, err := baseConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, ConfigDirName), nil
}

func TemplatesPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, TemplatesFileName), nil
}

func TemplatesDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, TemplatesDirName), nil
}

func ProjectsFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ProjectsFileName), nil
}

func CaddyfilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, CaddyfileName), nil
}

func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}

func ProxyLogPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ProxyLogName), nil
}

func ProxyErrPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ProxyErrName), nil
}

func ProxyPlistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, launchAgentsDir, ProxyLabel+".plist"), nil
}

func baseConfigDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}
