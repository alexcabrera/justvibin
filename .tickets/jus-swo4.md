---
id: jus-swo4
status: closed
deps: [jus-cys3]
links: []
created: 2026-02-12T19:59:47Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-2e3m
tags: [templates, django]
---
# Add justvibin.toml to django-hypermedia template

Add justvibin.toml manifest to the justvibin-with-django-and-hypermedia template.

## Repository

alexcabrera/justvibin-with-django-and-hypermedia

## Manifest Content

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
    ".DS_Store",
    ".venv",
    "db.sqlite3",
    "staticfiles",
    "node_modules"
]
setup = "./setup.sh"
setup_interactive = true

[serve]
type = "command"
dev = "./start.sh --dev"
prod = "./start.sh"
port_env = "PORT"
default_port = 8000

[serve.static]
# Fallback not typically used, but define for completeness
root = "staticfiles"
```

## Additional Changes

1. Update start.sh to respect PORT environment variable
   - Already uses PORT, but verify

2. Update AGENTS.md with new justvibin commands:
   - justvibin start (instead of ./start.sh --dev)
   - justvibin open
   - justvibin tunnel

3. Update README.md:
   - Document justvibin integration
   - Keep standalone instructions as alternative

4. Add .justvibin to .gitignore (project marker is per-instance)

## Acceptance Criteria

- [ ] justvibin.toml created with correct schema
- [ ] start.sh respects PORT env var
- [ ] AGENTS.md updated with justvibin commands
- [ ] README.md documents justvibin usage
- [ ] .justvibin added to .gitignore

