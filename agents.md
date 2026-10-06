# Project: nova-pkg

I'm building **nova-pkg**, a universal package installer for Arch Linux written in Go. One CLI wraps multiple package sources (pacman, AUR, Flatpak, Snap, AppImage) behind a common interface. Please continue the project in a way that matches the decisions below.

## Environment
- OS: Arch Linux (Omarchy), Go 1.27 installed via pacman at /usr/bin/go
- Module path: `github.com/J-V-Studios/nova-pkg`
- Repo is on GitHub under the J-V-Studios org, default branch `main`
- `GOBIN` is set to `~/.local/bin`, so `go install ./cmd/nova-pkg` puts the binary on my PATH
- I'm new to Go and to CI, so explain non-obvious choices briefly and keep changes small

## Layout
```
nova-pkg/
├── cmd/nova-pkg/main.go          # CLI entry point
├── internal/backend/
│   ├── backend.go                # Backend interface + Package struct
│   ├── pacman.go                 # pacman backend (done)
│   └── pacman_test.go            # tests parseSearch
├── go.mod
└── .github/workflows/ci.yml      # Arch container: go vet, go build, go test
```

## Architecture
Every package source implements this interface:

```go
type Package struct {
    Name, Version, Description, Source string
}

type Backend interface {
    Name() string
    Available() bool                      // exec.LookPath check for the tool
    Search(query string) ([]Package, error)
    Install(pkg string) error
}
```

- The CLI loops over backends, skips any where `Available()` is false, and aggregates results.
- Parsing of external tool output lives in its own pure function (e.g. `parseSearch`) so it can be unit tested without running the real tool.
- Backends shell out with `os/exec` unless a native API is better (the AUR uses its HTTP RPC API instead of calling yay/paru).

## Current state
- Done: interface, pacman backend (`search`, `install`), a unit test for `parseSearch`, a simple `os.Args` CLI with `search` and `install`, CI workflow
- `go vet`, `go build`, and `go test` all pass; code is gofmt-formatted
- CLI usage: `nova-pkg search <term>`, `nova-pkg install <name>`

## Conventions
- Run `gofmt -w .` before finishing any change; CI fails on unformatted code
- Tests live next to the code (`foo.go` and `foo_test.go`), and must pass with `go test ./...`
- Never require root to run tests. Anything needing sudo or the network goes behind a build tag, e.g. `//go:build integration && flatpak`
- Tests must skip cleanly if a tool is missing (`exec.LookPath` + `t.Skip`)
- Use the standard library where reasonable; add dependencies only when they clearly pay off (cobra for the CLI later is fine)
- Wrap errors with context (`fmt.Errorf("pacman search: %w", err)`)
- Don't add features beyond the task I ask for

## Roadmap (do in this order, one step at a time)
1. **AUR backend**: search via `https://aur.archlinux.org/rpc/?v=5&type=search&arg=<q>` using `net/http` + `encoding/json`; install by cloning the AUR git repo and running `makepkg -si` (never as root)
2. **Flatpak backend**: wrap `flatpak search` and `flatpak install`
3. **Parallel search**: query all available backends concurrently with goroutines, then merge and sort the results
4. **More commands**: `remove`, `update`, `list`
5. **Switch the CLI to cobra** once there are more than a few commands
6. **Config file** (`~/.config/nova-pkg/config.toml`) for backend priority, such as preferring pacman over flatpak
7. Snap and AppImage backends

## How I want you to work
- Before changing anything, read the existing files so new code matches the existing style
- Make the smallest change that completes the task, then tell me which files you changed
- After each change, run `gofmt -w .`, `go vet ./...`, and `go test ./...`, and report the results
- Add or update a unit test for any new parsing or logic
- If a task has several reasonable approaches, briefly list them and recommend one before writing code

## First task
Implement roadmap step 1: the AUR backend in `internal/backend/aur.go`, with a test for the JSON parsing using a sample response (no network in unit tests). Register it in `cmd/nova-pkg/main.go` alongside pacman.