---
id: jus-k78p
status: open
deps: [jus-2e3m]
links: []
created: 2026-02-12T19:53:53Z
type: epic
priority: 1
assignee: Alex Cabrera
parent: jus-nq0k
tags: [github, homebrew, distribution]
---
# GitHub & Distribution

Handle all GitHub repository management and distribution updates for the unified justvibin architecture.

## Repository Changes

### alexcabrera/justvibin
- Becomes the unified CLI (scaffolding + serving)
- Major version bump to 1.0.0
- New release with updated functionality

### alexcabrera/srv
- Archive repository
- Update README with deprecation notice and redirect to justvibin
- Set repository description to indicate archived status

### alexcabrera/homebrew-tap
- Update justvibin formula for new version
- Remove srv formula (or redirect)
- Update SHA256 checksums

## Release Process

1. Complete all code changes in justvibin
2. Tag release v1.0.0
3. Update Homebrew formula
4. Archive srv repository
5. Announce in README/releases

## Default Templates

Consider shipping default templates with justvibin:
- Option A: Bundle templates in the CLI repo
- Option B: Auto-install official templates on first run
- Option C: Just document how to install (current approach)

Recommend Option B: On 'justvibin install' or first 'justvibin new', offer to install official templates.

## Acceptance Criteria

- [ ] justvibin v1.0.0 released on GitHub
- [ ] Homebrew formula updated and working
- [ ] srv repository archived with redirect notice
- [ ] Installation works: brew install alexcabrera/tap/justvibin
- [ ] README documents the unified architecture

