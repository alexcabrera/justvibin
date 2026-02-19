# justvibin - Agent Guidelines

This document provides guidelines for AI agents working on the justvibin codebase.

---

## Project Overview

justvibin is a CLI tool written in Go that:
1. Scaffolds web projects from git-based templates
2. Manages a local HTTPS proxy (Caddy) for `*.localhost` routing
3. Handles project lifecycle (start, stop, open, tunnel)

## Architecture

```
cmd/justvibin/           # CLI commands (Cobra)
├── main.go              # Entry point
├── root.go              # Root command + global flags
├── cmd_*.go             # Command handlers
├── *.go                 # Command-specific logic
└── *_test.go            # Tests

internal/
├── config/              # Configuration paths and templates
├── exec/                # Process execution helpers
├── fsutil/              # Filesystem utilities
├── logging/             # Styled logging (Charm)
├── manifest/            # justvibin.toml parsing/validation
├── proxy/               # Caddy proxy + launchd integration
├── registry/            # Project registry (projects.json)
├── serve/               # Static file server + status
├── ui/                  # Charm UI components
└── version/             # Version info
```

## Key Files

| File | Purpose |
|------|---------|
| `cmd/justvibin/root.go` | Root command, global flags |
| `cmd/justvibin/cmd_new.go` | `justvibin new` command |
| `cmd/justvibin/cmd_start.go` | `justvibin start` command |
| `cmd/justvibin/cmd_stop.go` | `justvibin stop` command |
| `internal/manifest/manifest.go` | TOML parsing for justvibin.toml |
| `internal/registry/registry.go` | Project registry management |
| `internal/proxy/caddy.go` | Caddyfile generation |
| `internal/proxy/launchd.go` | macOS launchd service |
| `internal/config/paths.go` | Config directory paths |

## Adding Commands

1. Create `cmd/justvibin/cmd_<name>.go`:

```go
package main

import (
    "github.com/alexcabrera/justvibin/internal/logging"
    "github.com/alexcabrera/justvibin/internal/ui"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:     "mycommand [args]",
    Short:   "Brief description",
    Long:    "Detailed description for help output.",
    Example: "justvibin mycommand --flag value",
    Args:    cobra.MaximumNArgs(1),
    RunE:    runMyCmd,
}

func init() {
    rootCmd.AddCommand(myCmd)
    myCmd.Flags().StringP("flag", "f", "default", "Flag description")
}

func runMyCmd(cmd *cobra.Command, args []string) error {
    output := getOutputSettings(cmd)
    console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
    logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
    logger.SetSilent(output.Quiet)
    logger.SetVerbose(output.Verbose)

    // Implementation here
    return nil
}
```

2. Add tests in `cmd/justvibin/cmd_<name>_test.go`

## Template System

### Template Structure

Templates are git repositories containing:
- `justvibin.toml` — Manifest describing scaffold/serve behavior
- Project files to copy

### Manifest Schema

```go
type Manifest struct {
    Template Template  // name, description, version, author, url
    Scaffold Scaffold  // exclude, setup, setup_interactive
    Serve    Serve     // type, dev, prod, port_env, default_port
    Project  Project   // marker_fields
}
```

### Serve Types

1. **static** — Built-in Go HTTP server serves files
2. **command** — Runs external command (e.g., `./start.sh`)

## Registry System

Projects are tracked in `~/.config/justvibin/projects.json`:

```json
{
  "projects": [
    {
      "name": "myproject",
      "path": "/path/to/myproject",
      "port": 8001,
      "created": "2026-02-18T12:00:00Z"
    }
  ]
}
```

The `.justvibin` marker file in each project contains:
```
template = django-hypermedia
port = 8001
created = 2026-02-18T12:00:00Z
```

## Proxy System

- **Caddy** handles HTTPS termination and routing
- **launchd** manages Caddy as a system service on macOS
- Caddyfile regenerated when projects start/stop

## Common Tasks

### Run Tests

```bash
go test ./...
```

### Build

```bash
go build -o justvibin ./cmd/justvibin
```

### Add Dependency

```bash
go get <package>
go mod tidy
```

## Testing Patterns

Tests use helper functions from `root_test_helpers.go`:

```go
func TestMyCommand(t *testing.T) {
    tmpDir := t.TempDir()
    
    // Create test fixtures
    // ...
    
    // Run command
    rootCmd.SetArgs([]string{"mycommand", "--flag", "value"})
    err := rootCmd.Execute()
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

## Critical Rules

1. **Use Cobra patterns** — All commands use `spf13/cobra` with `RunE`
2. **Use Fang for execution** — `fang.Execute()` handles context
3. **Respect output flags** — Use `getOutputSettings()` for quiet/verbose/json
4. **Use internal packages** — Don't duplicate logging, UI, fsutil code
5. **Run tests** — `go test ./...` must pass before changes are complete
6. **Handle errors** — Return errors, don't panic; use logger for user messages

## Output Formatting

Commands support multiple output modes:
- **Default** — Styled terminal output (Charm)
- **--quiet** — Minimal output
- **--verbose** — Extra debug info
- **--json** — Machine-readable JSON
- **--no-color** — Plain text (for pipes/CI)

Use `getOutputSettings(cmd)` to get the current mode.

## Dependencies

| Package | Purpose |
|---------|---------|
| `spf13/cobra` | Command-line structure |
| `charmbracelet/fang` | Cobra execution wrapper |
| `charmbracelet/lipgloss` | Terminal styling |
| `charmbracelet/log` | Styled logging |
| `charmbracelet/huh` | Interactive prompts |
| `pelletier/go-toml/v2` | TOML parsing (not used, custom parser) |

## File Locations

| Path | Purpose |
|------|---------|
| `~/.config/justvibin/` | Config directory |
| `~/.config/justvibin/templates/` | Installed templates |
| `~/.config/justvibin/projects.json` | Project registry |
| `~/.config/justvibin/templates.toml` | Custom template definitions |
| `~/Library/LaunchAgents/land.charm.justvibin-proxy.plist` | launchd service |
| `/tmp/justvibin-proxy.Caddyfile` | Generated Caddyfile |
