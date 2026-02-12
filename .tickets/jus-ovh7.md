---
id: jus-ovh7
status: open
deps: [jus-jngy]
links: []
created: 2026-02-12T19:55:57Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [commands, templates]
---
# Implement 'justvibin install' command

Implement the install command for adding template plugins from git repositories.

## Usage

```
justvibin install <git-url>              # Install template
justvibin install <git-url> --name foo   # Install with custom name
justvibin install --list-official        # Show official templates
```

## Implementation

```bash
cmd_install() {
    local url="$1"
    local name=""
    
    # Parse --name option
    # ...
    
    # Clone to temp directory first
    local tmp_dir
    tmp_dir=$(mktemp -d)
    git clone --depth 1 "$url" "$tmp_dir"
    
    # Validate manifest exists
    if [[ ! -f "$tmp_dir/justvibin.toml" ]]; then
        log_error "Template missing justvibin.toml manifest"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    # Parse name from manifest if not provided
    if [[ -z "$name" ]]; then
        name=$(parse_toml_value "$tmp_dir/justvibin.toml" "template.name")
    fi
    
    # Check if already installed
    local target="$JUSTVIBIN_TEMPLATES_DIR/$name"
    if [[ -d "$target" ]]; then
        if confirm "Template '$name' already installed. Update?"; then
            rm -rf "$target"
        else
            rm -rf "$tmp_dir"
            exit 0
        fi
    fi
    
    # Move to templates directory
    mv "$tmp_dir" "$target"
    
    # Store source URL for updates
    echo "$url" > "$target/.source"
    
    log_success "Installed template: $name"
}
```

## Official Templates

Maintain list of official templates:

```bash
OFFICIAL_TEMPLATES=(
    "https://github.com/alexcabrera/justvibin-with-django-and-hypermedia.git"
    "https://github.com/alexcabrera/justvibin-with-hypertext.git"
)
```

## Template Directory Structure

```
~/.config/justvibin/templates/
├── django-hypermedia/
│   ├── .source              # URL for updates
│   ├── justvibin.toml       # Manifest
│   └── ...                  # Template files
└── hypertext/
    ├── .source
    ├── justvibin.toml
    └── ...
```

## Acceptance Criteria

- [ ] Templates cloned from git URLs
- [ ] Manifest validated before install
- [ ] Name derived from manifest or --name flag
- [ ] Existing templates can be updated
- [ ] Source URL stored for later updates
- [ ] --list-official shows available templates

