---
id: jus-cdr3
status: open
deps: []
links: []
created: 2026-02-12T19:53:06Z
type: epic
priority: 1
assignee: Alex Cabrera
parent: jus-nq0k
tags: [infrastructure, caddy, proxy]
---
# Infrastructure: Caddy Proxy System

Implement the central Caddy proxy infrastructure from srv into justvibin. This provides automatic HTTPS for all projects via *.localhost domains.

## Components

1. **Central Proxy Service**
   - Caddy running as launchd service
   - Routes https://<project>.localhost to project ports
   - Auto-generates Caddyfile from project registry

2. **Project Registry**
   - JSON file tracking all projects (name, port, path, template)
   - Port assignment (starting from 3000)
   - Project discovery via .justvibin marker files

3. **Certificate Management**
   - Caddy's local_certs for automatic HTTPS
   - One-time CA trust setup (sudo caddy trust)

## Files

- ~/.config/justvibin/projects.json - Project registry
- ~/.config/justvibin/Caddyfile - Auto-generated proxy config
- ~/Library/LaunchAgents/land.charm.justvibin.proxy.plist - launchd service

## Commands Affected

- justvibin install (setup proxy on first run)
- justvibin proxy start|stop|status|restart
- justvibin start/stop (register/unregister with proxy)
- justvibin list (show project status from registry)

## Acceptance Criteria

- [ ] Central Caddy proxy runs as launchd service
- [ ] Projects accessible via https://<name>.localhost
- [ ] Caddyfile auto-generated from project registry
- [ ] Port assignment avoids conflicts
- [ ] Proxy survives system restart

