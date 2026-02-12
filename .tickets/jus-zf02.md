---
id: jus-zf02
status: open
deps: [jus-cys3]
links: []
created: 2026-02-12T19:59:57Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-2e3m
tags: [templates, static]
---
# Add justvibin.toml to hypertext template

Add justvibin.toml manifest to the justvibin-with-hypertext template.

## Repository

alexcabrera/justvibin-with-hypertext

## Manifest Content

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
# No setup script needed for static sites

[serve]
type = "static"
# No commands needed - justvibin serves directly via Caddy

[serve.static]
root = "."
```

## Changes to Make

1. **Remove serve.sh**
   - No longer needed, justvibin handles serving
   - Static sites served directly by Caddy

2. **Update AGENTS.md**
   - Replace serve.sh references with justvibin commands
   - Document: justvibin start, justvibin open, justvibin tunnel

3. **Update README.md**
   - Installation via justvibin
   - Development workflow with justvibin
   - Keep fallback instructions (python -m http.server)

4. **Add .justvibin to .gitignore**

## Standalone Fallback

For users without justvibin, document:
```bash
# Without justvibin
python3 -m http.server 8000
# Then visit http://localhost:8000
```

## Acceptance Criteria

- [ ] justvibin.toml created with static type
- [ ] serve.sh removed
- [ ] AGENTS.md updated with justvibin commands
- [ ] README.md documents justvibin and standalone usage
- [ ] .justvibin added to .gitignore

