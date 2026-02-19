---
id: jus-thwh
status: closed
deps: [jus-ohdb, jus-c6us]
links: []
created: 2026-02-12T20:35:32Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fgp8
tags: [testing, static]
---
# Test hypertext with justvibin

Test the hypertext template end-to-end with justvibin.

## Test Steps

```bash
# 1. Install template
justvibin install https://github.com/alexcabrera/justvibin-with-hypertext.git
justvibin templates  # Verify listed

# 2. Create project
cd /tmp
justvibin new testhtml --template hypertext
# Verify: No setup script runs (static)
# Verify: Files copied

# 3. Start server
cd testhtml
justvibin start
# Verify: https://testhtml.localhost loads
# Verify: index.html shows
# Verify: HTMX demo works (partials load)

# 4. Test commands
justvibin stop
justvibin start --bg
justvibin open

# 5. Cleanup
justvibin remove testhtml --files
```

## Acceptance Criteria

- [ ] Template installs correctly
- [ ] Project creates without setup
- [ ] Static serving works
- [ ] HTMX partials work
- [ ] HTTPS works
- [ ] Cleanup works

