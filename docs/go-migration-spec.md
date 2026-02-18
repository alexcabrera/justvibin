# Go Migration Spec: justvibin v0.1.0 (Bash parity)

This document captures the current behavior of the Bash CLI (`justvibin` v0.1.0) to guide the Go port. Source of truth: `justvibin` Bash script in this repo.

## Scope

- Commands: `new`, `templates`, `--help/-h`, `--version/-v`.
- Configuration: `~/.config/justvibin/templates.toml` overrides.
- Default templates and display names.
- Gum-enhanced UX when available, ANSI fallback otherwise.
- Project creation flow: clone/copy, `.git` removal, `git init`, `setup.sh` execution, next steps.

## Global Behavior

- Requires Bash 4.0+ for associative arrays. Exits with error if not present.
- `JUSTVIBIN_VERSION="0.1.0"`.
- `JUSTVIBIN_CONFIG_DIR="$HOME/.config/justvibin"`.
- Colors defined for ANSI fallback (RED/GREEN/YELLOW/BLUE/MAGENTA/CYAN/BOLD/NC).

### Gum and TTY detection

- `has_gum`: returns true if `gum` is available on `PATH`.
- `is_interactive`: true when both stdin and stdout are TTYs (`-t 0` and `-t 1`).
- `can_use_gum`: `has_gum && is_interactive`.
- Note: some commands use `has_gum` (not `can_use_gum`) and will call `gum` even in non-tty contexts.

## Configuration

### Defaults

**Template URLs**
- `django-hypermedia` → `https://github.com/alexcabrera/justvibin-with-django-and-hypermedia.git`
- `hypertext` → `https://github.com/alexcabrera/justvibin-with-hypertext.git`

**Template Display Names**
- `django-hypermedia` → `Django + Hypermedia (HTMX, TailwindCSS)`
- `hypertext` → `HTML + Hypertext (Vanilla HTML/CSS/JS with HTMX)`

### User Overrides

`load_user_templates()` reads `${JUSTVIBIN_CONFIG_DIR}/templates.toml` if present. Parsing is line-based:
- Splits on `=`; trims whitespace; strips `"`.
- Ignores empty keys/values, `[templates]` header, and comment lines (`#`).
- Any key/value is accepted, no schema validation.
- Overrides or adds entries in `TEMPLATE_URLS`.

## UI / Logging

### Header

`print_header()`:
- If `has_gum`: renders a stylized double-border card with version.
- Else: prints a Magenta ASCII box with `justvibin` and version.

### Logging

`log_info`, `log_success`, `log_warn`, `log_error`:
- Use `gum style` when `has_gum` is true.
- Else print with ANSI colors (error goes to stderr).

### Spinner

`spin "Message" cmd...`:
- If `can_use_gum`: uses `gum spin` with dot spinner.
- Else prints `Message...` and reports `done`/`failed` based on command exit code.

## Commands

### `--help` / `-h`

- Prints header then help.
- Help sections (with gum styling or ANSI fallback):
  - Usage
  - Options for new
  - Available Templates (hard-coded list order)
  - Configuration (shows `~/.config/justvibin/templates.toml`)

### `--version` / `-v`

- Prints `justvibin v0.1.0` and exits 0.

### `templates`

- Prints `Available Templates` header.
- Iterates over `TEMPLATE_NAMES` associative array (order not guaranteed in Bash).
- For each template:
  - Name
  - Display name
  - URL (from `TEMPLATE_URLS`, or `<not configured>` if missing)
- Prints guidance snippet for `templates.toml`:
  -
    ```toml
    [templates]
    my-template = "https://github.com/user/repo.git"
    ```

### `new`

#### Argument parsing

- Accepts:
  - `--local PATH`
  - `--template NAME`
  - First positional argument becomes `project_name`.
- Unknown options (`-*`) => error + exit 1.

#### Flow

1. Print header **before** prompts/validation.
2. If `project_name` missing:
   - `prompt_project_name()`
     - `gum input` (placeholder `my-project`) when `can_use_gum`.
     - `read -rp` otherwise.
   - Empty input => error + exit 1.
3. Validate name regex: `^[a-zA-Z][a-zA-Z0-9_-]*$`.
4. If directory exists: error + exit 1.
5. If `--local` not provided and default template `django-hypermedia` is in use and multiple templates exist, prompt:
   - `choose_template()`
     - `gum choose` when `can_use_gum`.
     - `select` menu otherwise.
   - Maps display name back to template key.
6. Log project + template.
7. Copy/clone:
   - Local: `cp -R PATH project_name` (error if path missing).
   - Remote: `git clone --depth 1 URL project_name` (error if URL missing/unknown).
   - Always `rm -rf project_name/.git` after copy/clone.
8. Initialize repo: `git -C project_name init`.
9. Setup:
   - If `setup.sh` exists and is executable: `./setup.sh`.
   - Else if `setup.sh` exists: `bash setup.sh`.
   - Else: warn `No setup.sh found in template`.
10. Print success + next steps (`cd project_name` and `./start.sh --dev`).

## Exit Codes

- `--help`, `--version`, and no-args path: exit 0.
- Unknown command: error + help + exit 1.
- Most validation and command failures in `new` exit 1.

## Go Port Parity Notes

### Gum → Charm replacements

The Go port should replace gum with Charm libraries while keeping UX equivalent:

- `gum style` → **lipgloss** for layout/color.
- `gum choose` / `gum input` → **huh** forms/selects.
- `gum spin` → **bubbles/spinner** + **bubbletea** runtime.
- Logging format/colors → **charm/log** (optionally styled with lipgloss).

### Parity expectations

- Keep headers, help sections, template list, and config guidance consistent with Bash output.
- Prompt behavior should be TTY-gated to avoid blocking in non-interactive contexts.
- The Bash `templates` ordering is not deterministic; the Go port should pick a deterministic order (defaults first) to match help output ordering.
- Bash does not provide JSON output; any JSON output in Go is additive and must be explicitly flagged.

## Open Questions / Decisions

- Decide deterministic ordering for template lists (defaults then overrides recommended).
- Decide how to handle gum-available-but-non-tty scenarios (Bash may call gum; Go should avoid non-tty prompts).
- Decide whether to preserve the Bash `cp -R` semantics exactly or use Go-native copy with similar permissions.
