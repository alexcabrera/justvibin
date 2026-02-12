# justvibin

A CLI for scaffolding web application projects from curated templates.

## Installation

```bash
# Download the script
curl -o justvibin https://raw.githubusercontent.com/alexcabrera/justvibin/main/justvibin
chmod +x justvibin

# Optionally move to PATH
mv justvibin /usr/local/bin/
```

## Usage

```bash
# Create a new project (interactive)
justvibin new myproject

# Use a local template (for development)
justvibin new myproject --local /path/to/template

# List available templates
justvibin templates

# Show help
justvibin --help

# Show version
justvibin --version
```

## Available Templates

- **django-hypermedia** - Django + HTMX + TailwindCSS starter

## Custom Templates

Create `~/.config/justvibin/templates.toml` to add custom templates:

```toml
[templates]
my-template = "https://github.com/user/repo.git"
```

## Requirements

- bash 4.0+
- git
- [gum](https://github.com/charmbracelet/gum) (optional, for prettier UI)

## License

MIT
