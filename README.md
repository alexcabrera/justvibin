# justvibin

A CLI for scaffolding web application projects from curated templates.

## Installation

```bash
# Install the Go binary
# (replace with your GOPATH/bin or preferred install path)
go install github.com/alexcabrera/justvibin/cmd/justvibin@latest
```

## Building from source

```bash
go build ./...
```

## Usage

```bash
# Create a new project (interactive)
justvibin new myproject

# Use a local template (for development)
justvibin new myproject --local /path/to/template

# List available templates
justvibin templates
justvibin templates --json

# Show help
justvibin --help

# Show version
justvibin --version
```

## Available Templates

- **django-hypermedia** - Django + HTMX + TailwindCSS starter
- **hypertext** - Vanilla HTML/CSS/JS + HTMX

## Custom Templates

Create `~/.config/justvibin/templates.toml` to add custom templates:

```toml
[templates]
my-template = "https://github.com/user/repo.git"
```

## Template Manifest (justvibin.toml)

Templates include a `justvibin.toml` manifest describing scaffold and serve behavior.

```toml
[template]
name = "my-template"
description = "My awesome template"
version = "1.0.0"
author = "Alex Cabrera"
url = "https://github.com/user/repo.git"

[scaffold]
exclude = [".git", ".DS_Store"]
setup = "./setup.sh"
setup_interactive = true

[serve]
type = "command"
dev = "./start.sh --dev"
prod = "./start.sh"
port_env = "PORT"
default_port = 8000

[serve.static]
root = "."
extensions = [".html"]

[project]
marker_fields = ["template", "port", "created"]
```

Minimal manifest:

```toml
[template]
name = "my-template"
description = "My awesome template"

[serve]
type = "static"
```

Validation rules:

- `template.name` is required and must match `[a-z0-9-]+`
- `template.description` is required
- `serve.type` is required and must be `static` or `command`
- For `command` templates, `serve.dev` or `serve.prod` is required

## Requirements

- Go 1.21+
- git
- Charm UI deps (lipgloss/log/huh/bubbles) are bundled in the Go build

## License

MIT
