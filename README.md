**Someday CLI**

A minimal Go command line tool scaffolded with module path `github.com/kuniyoshi/someday`.

Usage
- Build: `go build -o bin/someday ./cmd/someday`
- Run: `go run ./cmd/someday help`
- Version: `go run ./cmd/someday version`

Notes
- The version string can be overridden at build time:
  `go build -ldflags "-X main.version=1.0.0" -o bin/someday ./cmd/someday`

