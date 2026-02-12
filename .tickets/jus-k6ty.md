---
id: jus-k6ty
status: open
deps: [jus-8hqk]
links: []
created: 2026-02-12T20:36:52Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ju43
tags: [github, deprecation]
---
# Update srv repository with deprecation notice

Update srv repository README with deprecation notice.

## New README.md

```markdown
# \u26a0\ufe0f srv has been merged into justvibin

This project's functionality is now part of [justvibin](https://github.com/alexcabrera/justvibin).

## Migration

\`\`\`bash
# Uninstall srv
brew uninstall srv

# Install justvibin
brew tap alexcabrera/tap
brew install justvibin

# Run setup (migrates your projects)
justvibin setup
\`\`\`

## Command Mapping

| srv | justvibin |
|-----|-----------|
| srv new | justvibin new |
| srv start | justvibin start |
| srv stop | justvibin stop |
| srv list | justvibin list |
| srv open | justvibin open |
| srv tunnel | justvibin tunnel |
| srv proxy | justvibin proxy |
| srv register | justvibin register |
| srv remove | justvibin remove |
| srv sync | justvibin sync |

## What's New in justvibin

- Template plugin system (install templates from git)
- Framework-aware serving (Django, Node, etc.)
- Better scaffolding

## Questions?

[Open an issue on justvibin](https://github.com/alexcabrera/justvibin/issues)
\`\`\`
```

## Acceptance Criteria

- [ ] README replaced with deprecation notice
- [ ] Migration steps documented
- [ ] Command mapping shown
- [ ] Links to justvibin work

