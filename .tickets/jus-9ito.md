---
id: jus-9ito
status: open
deps: [jus-jjl8]
links: []
created: 2026-02-12T20:02:18Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-k78p
tags: [documentation, ai]
---
# Add justvibin AGENTS.md

Create AGENTS.md for AI agents working on or with justvibin.

## File Location

/Users/acabrera/Code/justvibin/justvibin/AGENTS.md

## Content

```markdown
# justvibin - AI Agent Guidelines

## Project Overview

justvibin is a bash CLI for scaffolding and serving web projects.
It combines project templating with local HTTPS development servers.

## Architecture

### Core Components

1. **Template Plugin System**
   - Templates installed to ~/.config/justvibin/templates/
   - Each template has justvibin.toml manifest
   - Scaffolding copies template to new project

2. **Project Registry**
   - ~/.config/justvibin/projects.json tracks projects
   - Each project has .justvibin marker file
   - Port assignment managed centrally

3. **Caddy Proxy**
   - Central reverse proxy for HTTPS
   - Routes *.localhost to project ports
   - Runs as launchd service

### Key Files

| File | Purpose |
|------|---------|
| justvibin | Main CLI script (~1500 lines bash) |
| ~/.config/justvibin/projects.json | Project registry |
| ~/.config/justvibin/templates/ | Installed templates |
| ~/.config/justvibin/Caddyfile | Generated proxy config |
| .justvibin | Per-project marker (JSON) |
| justvibin.toml | Template manifest (TOML) |

## Working with justvibin

### Adding Commands

Commands follow pattern:
\`\`\`bash
cmd_<name>() {
    # Parse arguments
    # Do work
    # Log result
}
\`\`\`

Add to main() case statement.

### Dependencies

- bash 4.0+ (associative arrays)
- jq (JSON processing)
- caddy (web server)
- gum (optional, UI)

### Testing Changes

\`\`\`bash
# Test directly
./justvibin <command>

# Test with new project
./justvibin new testproj --template hypertext
cd testproj
./justvibin start
./justvibin stop
\`\`\`

## Creating Templates

Templates must include justvibin.toml:

\`\`\`toml
[template]
name = "my-template"
description = "Description here"

[serve]
type = "static"  # or "command"
\`\`\`

For command-based templates:
\`\`\`toml
[serve]
type = "command"
dev = "./start.sh --dev"
port_env = "PORT"
\`\`\`

## Common Tasks

### Add new serve type
1. Update start_server() in justvibin
2. Add case for new type
3. Document in README

### Add new command
1. Create cmd_<name>() function
2. Add to main() case statement
3. Add to cmd_help() output
4. Update shell completions

### Debug proxy issues
\`\`\`bash
justvibin proxy status
cat ~/.config/justvibin/proxy.log
caddy validate --config ~/.config/justvibin/Caddyfile
\`\`\`
\`\`\`

## Acceptance Criteria

- [ ] AGENTS.md created in justvibin repo
- [ ] Architecture documented
- [ ] Key files listed
- [ ] Common tasks explained
- [ ] Template creation documented

