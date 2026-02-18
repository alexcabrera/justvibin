---
id: jus-3cof
status: closed
deps: []
links: []
created: 2026-02-12T20:27:25Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-f3br
tags: [infrastructure, setup]
---
# Define constants and global variables

Define all global constants and variables at the top of the justvibin script.

## Variables to Define

```bash
# Version
JUSTVIBIN_VERSION="1.0.0"

# Configuration paths
JUSTVIBIN_CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/justvibin"
JUSTVIBIN_PROJECTS_FILE="${JUSTVIBIN_CONFIG_DIR}/projects.json"
JUSTVIBIN_TEMPLATES_DIR="${JUSTVIBIN_CONFIG_DIR}/templates"
JUSTVIBIN_CADDYFILE="${JUSTVIBIN_CONFIG_DIR}/Caddyfile"
JUSTVIBIN_PLIST="${HOME}/Library/LaunchAgents/land.charm.justvibin.proxy.plist"
JUSTVIBIN_LOG="${JUSTVIBIN_CONFIG_DIR}/proxy.log"
JUSTVIBIN_ERR="${JUSTVIBIN_CONFIG_DIR}/proxy.err"

# Defaults
JUSTVIBIN_BASE_PORT=3000
JUSTVIBIN_PROXY_LABEL="land.charm.justvibin.proxy"

# Official templates
JUSTVIBIN_OFFICIAL_TEMPLATES=(
    "https://github.com/alexcabrera/justvibin-with-django-and-hypermedia.git"
    "https://github.com/alexcabrera/justvibin-with-hypertext.git"
)
```

## Considerations

- Respect XDG_CONFIG_HOME if set
- Use consistent naming convention
- Group related variables together

## Acceptance Criteria

- [ ] All paths defined as variables
- [ ] XDG_CONFIG_HOME respected
- [ ] Official templates array defined
- [ ] Version variable matches release plan

