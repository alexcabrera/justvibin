---
id: jus-6d0x
status: closed
deps: [jus-rg8c]
links: []
created: 2026-02-12T19:53:31Z
type: epic
priority: 1
assignee: Alex Cabrera
parent: jus-nq0k
tags: [serving, runners, framework]
---
# Framework-Aware Serving System

Implement framework-aware serving that reads project's template manifest to determine how to start the development server.

## Serving Modes

### 1. Static Mode (type = "static")
- Use Caddy file-server directly
- No external process needed
- Fastest for pure HTML/CSS/JS

### 2. Command Mode (type = "command")
- Run template-defined command (e.g., ./start.sh --dev)
- Pass PORT via environment variable
- justvibin manages the process lifecycle

## Project Detection

When running 'justvibin start' in a project directory:

1. Look for .justvibin file (contains template name, port)
2. Find template manifest in ~/.config/justvibin/templates/<name>/justvibin.toml
3. Read serve configuration
4. Start appropriate server

## Process Management

- Track running processes in ~/.config/justvibin/running.json
- PID files in project directories (.justvibin.pid)
- Graceful shutdown on 'justvibin stop'
- Status tracking in 'justvibin list'

## Port Management

- Central proxy routes https://<name>.localhost to assigned port
- Port passed to command via serve.port_env environment variable
- For static: Caddy file-server listens on port directly

## Fallback Behavior

If template not found or manifest missing:
1. Check for common markers (manage.py, package.json, index.html)
2. Use sensible defaults
3. Warn user about missing manifest

## Acceptance Criteria

- [ ] Static projects served via Caddy file-server
- [ ] Command-based projects run template-defined commands
- [ ] PORT passed correctly to project commands
- [ ] Process lifecycle managed (start/stop)
- [ ] Running status tracked and displayed in list

