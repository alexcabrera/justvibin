---
id: jus-2e3m
status: closed
deps: [jus-6d0x]
links: []
created: 2026-02-12T19:53:42Z
type: epic
priority: 1
assignee: Alex Cabrera
parent: jus-nq0k
tags: [templates, migration]
---
# Template Migrations

Update existing template repositories to work with the new plugin architecture. Each template needs a justvibin.toml manifest and potentially updated scripts.

## Templates to Migrate

1. **justvibin-with-django-and-hypermedia**
   - Add justvibin.toml with command-based serving
   - Keep existing start.sh (works as-is)
   - Update AGENTS.md with new justvibin commands

2. **justvibin-with-hypertext**
   - Add justvibin.toml with static serving
   - Remove serve.sh (replaced by justvibin infrastructure)
   - Update AGENTS.md and README.md

## Manifest Requirements

Each template must have justvibin.toml in root with:
- Template metadata (name, description)
- Scaffold configuration (exclude patterns, setup script)
- Serve configuration (static vs command)

## Backward Compatibility

Templates should still work standalone:
- Django template: ./start.sh --dev still works
- Static template: can still use python -m http.server

The manifest just enables integration with justvibin's enhanced features (HTTPS, tunnels, etc.)

## Acceptance Criteria

- [ ] django-hypermedia template has justvibin.toml
- [ ] hypertext template has justvibin.toml
- [ ] Both templates work with new justvibin CLI
- [ ] Both templates still work standalone
- [ ] Documentation updated in both templates

