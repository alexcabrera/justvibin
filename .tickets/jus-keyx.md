---
id: jus-keyx
status: open
deps: [jus-q0si, jus-ae77]
links: []
created: 2026-02-12T19:56:52Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [commands, templates, refactor]
---
# Refactor 'justvibin new' for plugin architecture

Refactor the 'new' command to work with the template plugin system.

## Current Behavior

- Templates hardcoded in TEMPLATE_URLS associative array
- Clones directly from GitHub
- Runs setup.sh if present

## New Behavior

1. List installed templates from ~/.config/justvibin/templates/
2. Show picker if multiple templates (or --template flag)
3. Copy from local template directory (not git clone)
4. Apply scaffold.exclude patterns
5. Write .justvibin project marker
6. Register with project registry
7. Run scaffold.setup if defined

## Usage

```
justvibin new [name]                    # Interactive
justvibin new myproject                 # Use default/only template
justvibin new myproject -t hypertext    # Specify template
justvibin new myproject --template hypertext
```

## Implementation Changes

```bash
cmd_new() {
    local project_name=""
    local template=""
    
    # Parse arguments...
    
    # Get available templates
    local templates=()
    for dir in "$JUSTVIBIN_TEMPLATES_DIR"/*/; do
        templates+=("$(basename "$dir")")
    done
    
    if [[ ${#templates[@]} -eq 0 ]]; then
        log_error "No templates installed"
        log_info "Install one: justvibin install --list-official"
        exit 1
    fi
    
    # Select template
    if [[ -z "$template" ]]; then
        if [[ ${#templates[@]} -eq 1 ]]; then
            template="${templates[0]}"
        else
            template=$(choose_template)
        fi
    fi
    
    # Validate template exists
    local template_dir="$JUSTVIBIN_TEMPLATES_DIR/$template"
    if [[ ! -d "$template_dir" ]]; then
        log_error "Template '$template' not found"
        exit 1
    fi
    
    # Copy template (respecting exclude patterns)
    copy_template "$template_dir" "$project_name"
    
    # Write .justvibin marker
    local port=$(next_port)
    write_project_marker "$project_name" "$template" "$port"
    
    # Register project
    register_project "$project_name" "$port" "$(pwd)/$project_name" "$template"
    
    # Regenerate Caddyfile
    generate_caddyfile
    reload_proxy
    
    # Run setup script if defined
    local setup_cmd=$(parse_toml_value "$template_dir/justvibin.toml" "scaffold.setup")
    if [[ -n "$setup_cmd" ]]; then
        cd "$project_name"
        log_info "Running setup: $setup_cmd"
        eval "$setup_cmd"
    fi
    
    log_success "Created $project_name"
    log_info "URL: https://$project_name.localhost"
}
```

## Template Copying

```bash
copy_template() {
    local src="$1"
    local dest="$2"
    
    # Get exclude patterns
    local excludes=()
    # Parse scaffold.exclude array...
    # Always exclude: .git, justvibin.toml, .source
    
    # Use rsync or manual copy with exclusions
    rsync -a --exclude='.git' --exclude='justvibin.toml' --exclude='.source' \
        "$src/" "$dest/"
}
```

## Acceptance Criteria

- [ ] Templates read from installed plugins
- [ ] Picker shown when multiple templates available
- [ ] --template flag selects specific template
- [ ] scaffold.exclude patterns respected
- [ ] .justvibin marker written with template, port
- [ ] Project registered in registry
- [ ] scaffold.setup script runs if defined
- [ ] Error if no templates installed

