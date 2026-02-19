---
id: jus-asej
status: closed
deps: []
links: []
created: 2026-02-19T18:59:13Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-dvrr
---
# Add 'justvibin port' command to update project port

Add a new CLI command 'justvibin port [--set PORT]' that can query or update the port for the current project. When updating, it modifies the registry, regenerates Caddyfile, and reloads Caddy.

## Acceptance Criteria

- 'justvibin port' shows current port
- 'justvibin port --set 8001' updates port
- Registry (projects.json) is updated
- .justvibin marker file is updated
- Caddyfile is regenerated
- Caddy proxy is reloaded

