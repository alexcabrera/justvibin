# justvibin

A CLI for scaffolding web applications from curated templates with local HTTPS development.

## Features

- **Project Scaffolding** — Create new projects from official or custom templates
- **Local HTTPS Proxy** — Automatic `https://<project>.localhost` URLs via Caddy
- **Framework-Aware Serving** — Handles static sites and command-based servers
- **Template Plugin System** — Install templates from any git repository
- **Cloudflare Tunnels** — Share projects publicly with one command

## Installation

```bash
# Homebrew (macOS)
brew tap alexcabrera/tap
brew install justvibin

# From source
go install github.com/alexcabrera/justvibin/cmd/justvibin@latest
```

## Quick Start

```bash
# One-time setup (installs Caddy proxy)
justvibin setup

# Create a new Django project
justvibin new myapp --template django-hypermedia

# Start the server
cd myapp
justvibin start

# Open in browser (https://myapp.localhost)
justvibin open
```

## Commands

| Command | Description |
|---------|-------------|
| `justvibin new <name>` | Create a new project from a template |
| `justvibin start` | Start the project server |
| `justvibin stop` | Stop the running server |
| `justvibin open` | Open project in browser |
| `justvibin list` | List all registered projects |
| `justvibin templates` | List installed templates |
| `justvibin install <url>` | Install a template from git URL |
| `justvibin uninstall <name>` | Remove an installed template |
| `justvibin update <name>` | Update a template from its source |
| `justvibin tunnel` | Share project via Cloudflare tunnel |
| `justvibin proxy start` | Start the HTTPS proxy service |
| `justvibin proxy stop` | Stop the proxy service |
| `justvibin setup` | First-time setup wizard |
| `justvibin register` | Register existing directory as project |
| `justvibin remove <name>` | Remove project from registry |
| `justvibin sync` | Rebuild registry by scanning for projects |

### Global Flags

| Flag | Description |
|------|-------------|
| `--quiet, -q` | Suppress non-essential output |
| `--verbose, -v` | Enable verbose output |
| `--json` | Output in JSON format (where supported) |
| `--no-color` | Disable colored output |

## Official Templates

| Template | Description |
|----------|-------------|
| `django-hypermedia` | Django + HTMX + TailwindCSS full-stack application |
| `hypertext` | Vanilla HTML/CSS/JS with HTMX — simple static site |

Install official templates:

```bash
justvibin install --list-official
justvibin install https://github.com/alexcabrera/justvibin-with-django-and-hypermedia.git
justvibin install https://github.com/alexcabrera/justvibin-with-hypertext.git
```

## Custom Templates

Add custom templates via `~/.config/justvibin/templates.toml`:

```toml
[templates]
my-template = "https://github.com/user/my-template.git"
```

## Creating Templates

Templates are git repositories with a `justvibin.toml` manifest:

```toml
[template]
name = "my-template"
description = "My awesome template"
version = "1.0.0"
author = "Your Name"
url = "https://github.com/user/my-template.git"

[scaffold]
exclude = [".git", "justvibin.toml", ".DS_Store"]
setup = "./setup.sh"           # Optional setup script
setup_interactive = true       # Requires TTY

[serve]
type = "command"               # "command" or "static"
dev = "./start.sh --dev"
prod = "./start.sh"
port_env = "PORT"              # Environment variable for port
default_port = 8000
```

For static sites:

```toml
[template]
name = "static-site"
description = "Simple static website"

[scaffold]
exclude = [".git", "justvibin.toml"]

[serve]
type = "static"

[serve.static]
root = "."
```

### Validation Rules

- `template.name` — Required, must match `[a-z0-9-]+`
- `template.description` — Required
- `serve.type` — Required, must be `static` or `command`
- For `command` type: `serve.dev` or `serve.prod` required

## How It Works

1. **Project Creation**: `justvibin new` clones a template, excludes specified files, and runs the setup script
2. **Registration**: Projects are registered in `~/.config/justvibin/projects.json` with their port assignments
3. **Serving**: `justvibin start` launches the server (static or command-based) and registers with the proxy
4. **Proxy**: Caddy runs as a launchd service, routing `*.localhost` to project ports with automatic HTTPS
5. **Tunnels**: `justvibin tunnel` uses Cloudflare's quick tunnel for temporary public URLs

## Requirements

- **macOS** (Linux support planned)
- **Go 1.21+** (for building from source)
- **git** (for cloning templates)
- **Caddy** (installed automatically via `justvibin setup`)
- **cloudflared** (optional, for tunnels)

## License

MIT
