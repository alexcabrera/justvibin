---
id: jus-ocjp
status: closed
deps: [jus-keyx, jus-wz9o]
links: []
created: 2026-02-12T19:57:39Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-6d0x
tags: [commands, serving]
---
# Implement 'justvibin start' command

Implement the start command for serving projects based on their template configuration.

## Usage

```
justvibin start              # Start current project (dev mode)
justvibin start --prod       # Start in production mode
justvibin start <name>       # Start specific project by name
```

## Implementation

```bash
cmd_start() {
    local project_name=""
    local prod_mode=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --prod) prod_mode=true; shift ;;
            *) project_name="$1"; shift ;;
        esac
    done
    
    # Get project info
    if [[ -z "$project_name" ]]; then
        # Current directory
        if ! is_justvibin_project; then
            log_error "Not a justvibin project directory"
            log_info "Run 'justvibin new' or 'justvibin register' first"
            exit 1
        fi
        project_name=$(get_current_project_name)
    fi
    
    local template=$(get_current_project_template)
    local port=$(get_current_project_port)
    local template_dir="$JUSTVIBIN_TEMPLATES_DIR/$template"
    
    # Get serve configuration
    local serve_type="static"
    if [[ -f "$template_dir/justvibin.toml" ]]; then
        serve_type=$(parse_toml_value "$template_dir/justvibin.toml" "serve.type")
    fi
    
    # Check if already running
    if is_project_running "$project_name"; then
        log_warn "Project already running"
        log_info "URL: https://$project_name.localhost"
        return
    fi
    
    # Start based on type
    case "$serve_type" in
        static)
            start_static_server "$port"
            ;;
        command)
            start_command_server "$template_dir" "$port" "$prod_mode"
            ;;
        *)
            log_error "Unknown serve type: $serve_type"
            exit 1
            ;;
    esac
    
    # Track running process
    track_running_process "$project_name" $!
    
    log_success "Started: https://$project_name.localhost"
}

start_static_server() {
    local port="$1"
    local root="."
    
    # Start Caddy file-server in background
    caddy file-server --listen ":$port" --root "$root" &
    echo $! > .justvibin.pid
}

start_command_server() {
    local template_dir="$1"
    local port="$2"
    local prod_mode="$3"
    
    local cmd
    if $prod_mode; then
        cmd=$(parse_toml_value "$template_dir/justvibin.toml" "serve.prod")
    else
        cmd=$(parse_toml_value "$template_dir/justvibin.toml" "serve.dev")
    fi
    
    local port_env=$(parse_toml_value "$template_dir/justvibin.toml" "serve.port_env")
    port_env="${port_env:-PORT}"
    
    # Export port and run command
    export "$port_env"="$port"
    eval "$cmd" &
    echo $! > .justvibin.pid
}
```

## Process Tracking

Store PID in:
- Project directory: .justvibin.pid
- Central registry: ~/.config/justvibin/running.json

## Foreground vs Background

For development, consider running in foreground so user sees output.
Use --background flag for background execution.

Or always background but tail logs with 'justvibin logs'.

## Acceptance Criteria

- [ ] Static projects served via Caddy
- [ ] Command projects run template command
- [ ] PORT passed via configured env var
- [ ] Already-running check works
- [ ] PID tracked for later stop
- [ ] URL displayed on start

