---
id: jus-q0si
status: open
deps: [jus-jngy]
links: []
created: 2026-02-12T19:57:07Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [templates, project]
---
# Implement .justvibin project marker file

Implement the .justvibin project marker file that identifies justvibin-managed projects.

## File Format

JSON format for easy parsing:

```json
{
  "name": "myproject",
  "template": "django-hypermedia",
  "port": 3000,
  "created": "2024-01-15T10:30:00Z"
}
```

## Functions

```bash
write_project_marker() {
    local project_dir="$1"
    local name="$2"
    local template="$3"
    local port="$4"
    
    local created=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    cat > "$project_dir/.justvibin" << EOF
{
  "name": "$name",
  "template": "$template",
  "port": $port,
  "created": "$created"
}
EOF
}

read_project_marker() {
    local project_dir="${1:-.}"
    local marker="$project_dir/.justvibin"
    
    if [[ ! -f "$marker" ]]; then
        return 1
    fi
    
    cat "$marker"
}

get_current_project_name() {
    jq -r '.name' .justvibin 2>/dev/null
}

get_current_project_port() {
    jq -r '.port' .justvibin 2>/dev/null
}

get_current_project_template() {
    jq -r '.template' .justvibin 2>/dev/null
}

is_justvibin_project() {
    [[ -f ".justvibin" ]]
}
```

## Migration from .srv

If .srv file exists (from old srv tool):
1. Read name and port from .srv
2. Create .justvibin with template="unknown"
3. Optionally remove .srv

## .gitignore

Template should include .justvibin in .gitignore? Or commit it?

Arguments for committing:
- Project settings are part of project
- Team members get same config

Arguments for ignoring:
- Port might differ per machine
- Personal preference

Recommend: Commit it, but allow port override via environment.

## Acceptance Criteria

- [ ] .justvibin marker created on 'new'
- [ ] Contains name, template, port, created
- [ ] Helper functions read marker values
- [ ] is_justvibin_project() detects projects
- [ ] Migration from .srv works

