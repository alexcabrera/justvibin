---
id: jus-lz4t
status: closed
deps: [jus-bx3x]
links: []
created: 2026-02-12T20:34:52Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-swo4
tags: [templates, django]
---
# Update django template .gitignore

Update the Django template's .gitignore to include justvibin files.

## Add to .gitignore

```
# justvibin
.justvibin.pid
```

## Note

The .justvibin marker file should be committed (it defines project config).
Only the .justvibin.pid file should be ignored.

## Acceptance Criteria

- [ ] .justvibin.pid added to .gitignore
- [ ] .justvibin NOT ignored (should be committed)

