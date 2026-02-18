---
id: jus-ae77
status: closed
deps: [jus-ovh7]
links: []
created: 2026-02-12T19:56:10Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [commands, templates]
---
# Implement 'justvibin templates' command

Implement the templates command for listing installed template plugins.

## Usage

```
justvibin templates          # List installed templates
justvibin templates --json   # Output as JSON
```

## Output Format

With gum:
```
Installed Templates

  django-hypermedia
    Django + HTMX + TailwindCSS
    Type: command
    Source: https://github.com/alexcabrera/justvibin-with-django-and-hypermedia.git

  hypertext
    Vanilla HTML/CSS/JS + HTMX
    Type: static
    Source: https://github.com/alexcabrera/justvibin-with-hypertext.git

To install more templates:
  justvibin install <git-url>
  justvibin install --list-official
```

## Implementation

```bash
cmd_templates() {
    local json_output=false
    [[ "${1:-}" == "--json" ]] && json_output=true
    
    local templates_dir="$JUSTVIBIN_TEMPLATES_DIR"
    
    if [[ ! -d "$templates_dir" ]] || [[ -z "$(ls -A "$templates_dir" 2>/dev/null)" ]]; then
        log_info "No templates installed."
        log_info "Install with: justvibin install <git-url>"
        log_info "Or see official: justvibin install --list-official"
        return
    fi
    
    if $json_output; then
        # Output JSON array
        echo "["
        local first=true
        for dir in "$templates_dir"/*/; do
            # ... build JSON
        done
        echo "]"
    else
        # Pretty output
        for dir in "$templates_dir"/*/; do
            local name=$(basename "$dir")
            local manifest="$dir/justvibin.toml"
            local description=$(parse_toml_value "$manifest" "template.description")
            local serve_type=$(parse_toml_value "$manifest" "serve.type")
            local source=$(cat "$dir/.source" 2>/dev/null || echo "unknown")
            
            # Display with gum or fallback
        done
    fi
}
```

## Acceptance Criteria

- [ ] Lists all installed templates
- [ ] Shows name, description, type, source URL
- [ ] Handles no templates gracefully
- [ ] --json flag outputs valid JSON
- [ ] Suggests install commands if empty

