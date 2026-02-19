---
id: jus-psut
status: closed
deps: []
links: []
created: 2026-02-12T20:35:00Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-zf02
tags: [templates, static]
---
# Create justvibin.toml for hypertext

Create the justvibin.toml manifest file for the hypertext template.

## File: justvibin.toml

```toml
[template]
name = "hypertext"
description = "Vanilla HTML/CSS/JS with HTMX - simple static site"
version = "1.0.0"
author = "Alex Cabrera"
url = "https://github.com/alexcabrera/justvibin-with-hypertext"

[scaffold]
exclude = [
    ".git",
    "justvibin.toml",
    ".justvibin",
    ".DS_Store"
]

[serve]
type = "static"

[serve.static]
root = "."
```

## Location

/Users/acabrera/Code/justvibin/justvibin-with-hypertext/justvibin.toml

## Acceptance Criteria

- [ ] justvibin.toml created
- [ ] serve.type is static
- [ ] Validated correctly

