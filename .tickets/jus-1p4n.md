---
id: jus-1p4n
status: open
deps: [jus-jjl8, jus-fgp8]
links: []
created: 2026-02-12T20:01:23Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-k78p
tags: [github, release]
---
# Create justvibin v1.0.0 release

Create the v1.0.0 release of unified justvibin on GitHub.

## Prerequisites

All these tickets must be complete:
- Infrastructure (proxy, registry, Caddyfile)
- Template plugin system
- Serving system
- Template migrations
- README update

## Release Checklist

### 1. Final Code Review

```bash
# Verify all commands work
justvibin --help
justvibin --version  # Should show 1.0.0
justvibin setup --check
```

### 2. Update Version

In justvibin script:
```bash
JUSTVIBIN_VERSION="1.0.0"
```

### 3. Create Git Tag

```bash
cd /Users/acabrera/Code/justvibin/justvibin
git add -A
git commit -m "Release v1.0.0: Unified CLI with template plugins and serving"
git tag -a v1.0.0 -m "justvibin v1.0.0

Major release: Unified CLI architecture

Features:
- Template plugin system (install from git)
- Automatic HTTPS via Caddy proxy
- Framework-aware serving (static, command)
- Cloudflare tunnel support
- Project registry and management

This release merges functionality from srv into justvibin.
"
git push origin main --tags
```

### 4. Create GitHub Release

Via GitHub UI or gh CLI:

```bash
gh release create v1.0.0 --title "justvibin v1.0.0" --notes "
## 🎉 Major Release: Unified CLI

justvibin is now a complete local development environment manager.

### Features

- **Template Plugins** - Install project templates from git repositories
- **Automatic HTTPS** - Every project gets \`https://name.localhost\`
- **Framework-Aware Serving** - Supports Django, Node, static sites, and more
- **Cloudflare Tunnels** - Share projects publicly with one command
- **Project Registry** - Track and manage all your projects

### Migration from srv

srv functionality has been merged into justvibin. Existing srv users:

\`\`\`bash
brew uninstall srv
brew upgrade justvibin
justvibin setup  # Migrates existing projects
\`\`\`

### Installation

\`\`\`bash
brew tap alexcabrera/tap
brew install justvibin
justvibin setup
\`\`\`

### Official Templates

- [django-hypermedia](https://github.com/alexcabrera/justvibin-with-django-and-hypermedia) - Django + HTMX + TailwindCSS
- [hypertext](https://github.com/alexcabrera/justvibin-with-hypertext) - Vanilla HTML/CSS/JS + HTMX
"
```

### 5. Get SHA256 for Homebrew

```bash
curl -sL https://github.com/alexcabrera/justvibin/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
```

## Acceptance Criteria

- [ ] Version updated to 1.0.0
- [ ] Git tag created and pushed
- [ ] GitHub release created with notes
- [ ] SHA256 calculated for Homebrew
- [ ] Release downloadable and working

