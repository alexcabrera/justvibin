---
id: jus-bx3x
status: open
deps: []
links: []
created: 2026-02-12T20:34:31Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-swo4
tags: [templates, django]
---
# Create justvibin.toml for django-hypermedia

Create the justvibin.toml manifest file for the Django template.

## File: justvibin.toml

```toml
[template]
name = "django-hypermedia"
description = "Django + HTMX + TailwindCSS full-stack web application"
version = "1.0.0"
author = "Alex Cabrera"
url = "https://github.com/alexcabrera/justvibin-with-django-and-hypermedia"

[scaffold]
exclude = [
    ".git",
    "justvibin.toml",
    ".justvibin",
    "__pycache__",
    "*.pyc",
    ".venv",
    "db.sqlite3",
    "staticfiles",
    "node_modules",
    ".read-only"
]
setup = "./setup.sh"
setup_interactive = true

[serve]
type = "command"
dev = "./start.sh --dev"
prod = "./start.sh"
port_env = "PORT"
default_port = 8000
```

## Location

/Users/acabrera/Code/justvibin/justvibin-with-django-and-hypermedia/justvibin.toml

## Acceptance Criteria

- [ ] justvibin.toml created
- [ ] All fields populated correctly
- [ ] Validated with validate_manifest()

