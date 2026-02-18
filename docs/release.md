# Release Guide (Go binary)

## Build matrix

```bash
GOOS=darwin  GOARCH=arm64 go build -o dist/justvibin-darwin-arm64 ./cmd/justvibin
GOOS=darwin  GOARCH=amd64 go build -o dist/justvibin-darwin-amd64 ./cmd/justvibin
GOOS=linux   GOARCH=arm64 go build -o dist/justvibin-linux-arm64 ./cmd/justvibin
GOOS=linux   GOARCH=amd64 go build -o dist/justvibin-linux-amd64 ./cmd/justvibin
```

## Versioning

- Source of truth: `internal/version.Version`.
- Tag releases as `vX.Y.Z` to match the version constant.

## Checksums

```bash
shasum -a 256 dist/justvibin-* > dist/justvibin-checksums.txt
```

## Optional Homebrew tap

- Create a tap repo and formula that downloads the desired artifact.
- Use `sha256` from the checksums file.
- Keep the formula version aligned with `internal/version.Version`.

No publishing is performed as part of this doc.
