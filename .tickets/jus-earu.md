---
id: jus-earu
status: closed
deps: [jus-asej]
links: []
created: 2026-02-19T18:59:08Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-dvrr
---
# Integrate justvibin CLI for Caddy updates in dev mode

When start.sh runs in dev mode and selects a new port, it should call 'justvibin update-port' (or similar) to update the project registry and regenerate the Caddyfile. This ensures the reverse proxy always points to the correct backend port.

## Acceptance Criteria

- start.sh detects if justvibin CLI is available
- Calls justvibin to update port in registry
- Caddy configuration is regenerated
- Caddy is reloaded to apply changes
- Falls back gracefully if justvibin not available

