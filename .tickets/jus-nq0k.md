---
id: jus-nq0k
status: open
deps: []
links: []
created: 2026-02-12T19:52:57Z
type: epic
priority: 0
assignee: Alex Cabrera
tags: [architecture, planning]
---
# Unified justvibin CLI Architecture

Merge srv and justvibin into a single unified CLI that handles both project scaffolding and serving. Templates become installable plugins that define how projects are created and served.

## Vision

justvibin becomes a complete local development environment manager:
- Install templates as plugins from git repositories
- Scaffold new projects from installed templates
- Serve projects with automatic HTTPS via Caddy (*.localhost)
- Framework-aware serving (static, Django, Node, etc.)
- Cloudflare tunnel support for sharing

## Current State

- **justvibin** (alexcabrera/justvibin): CLI for scaffolding from templates
- **srv** (alexcabrera/srv): Static site dev server with Caddy/HTTPS
- **justvibin-with-django-and-hypermedia**: Django template with complex start.sh
- **justvibin-with-hypertext**: Simple HTML/HTMX template

## Target State

Single justvibin CLI that:
1. Replaces both justvibin and srv
2. Has plugin architecture for templates
3. Serves any project type via template-defined runners
4. Maintains HTTPS routing via central Caddy proxy

## Key Commands (Target)

```
justvibin install <git-url>     # Install template plugin
justvibin uninstall <name>      # Remove template plugin
justvibin templates             # List installed templates
justvibin new [name]            # Scaffold project (picker if needed)
justvibin start                 # Serve current project
justvibin stop                  # Stop server
justvibin open                  # Open in browser
justvibin list                  # Show all projects
justvibin tunnel                # Cloudflare tunnel
justvibin proxy start|stop      # Manage central proxy
```

## Affected Repositories

- alexcabrera/justvibin (primary changes)
- alexcabrera/srv (deprecated, merged)
- alexcabrera/justvibin-with-django-and-hypermedia (template manifest)
- alexcabrera/justvibin-with-hypertext (template manifest)
- alexcabrera/homebrew-tap (formula updates)

## Acceptance Criteria

- [ ] Single justvibin CLI replaces both justvibin and srv
- [ ] Templates installable as plugins from git URLs
- [ ] Projects served with HTTPS via *.localhost
- [ ] Framework-aware serving via template manifests
- [ ] All existing functionality preserved
- [ ] Homebrew formula updated
- [ ] srv repository archived with redirect

