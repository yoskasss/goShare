# Development

Guide for contributing to GoShare development.

## Prerequisites

- Go 1.25+
- Git
- Make
- golangci-lint (optional, for linting)

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Project Structure

```
GoShare/
├── cmd/goshare/
│   └── main.go              # Entry point
├── internal/
│   ├── cli/
│   │   └── app.go           # CLI application logic
│   ├── project/
│   │   ├── manager.go       # Project management
│   │   └── manager_test.go  # Tests
│   ├── publish/
│   │   └── publisher.go     # Cloudflare publishing
│   ├── share/
│   │   ├── manager.go       # Share code management
│   │   └── manager_test.go  # Tests
│   └── ziputil/
│       ├── manager.go       # ZIP utilities
│       └── manager_test.go  # Tests
├── wiki/                    # Documentation
├── docs/legacy/             # Legacy web assets
├── bin/                     # Build output (gitignored)
├── go.mod / go.sum
├── Makefile
├── README.md
└── .gitignore
```

## Getting Started

### Clone and Build

```bash
git clone https://github.com/yourusername/GoShare.git
cd GoShare

# Install dependencies
go mod tidy

# Build
make build

# Run tests
make test

# Run in development mode
make dev ARGS="create -project test"
```

### Development Workflow

```bash
# 1. Create feature branch
git checkout -b feature/awesome-feature

# 2. Make changes
# Edit code...

# 3. Format and vet
make fmt
make vet

# 4. Run tests
make test

# 5. Build and test manually
make build
./bin/goshare help

# 6. Commit
git add .
git commit -m "feat: add awesome feature"

# 7. Push and create PR
git push origin feature/awesome-feature
```

## Makefile Targets

```bash
make help           # Show all targets
make build          # Build for current platform
make build-all      # Build for all platforms
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make fmt            # Format code (go fmt)
make vet            # Static analysis (go vet)
make lint           # Run golangci-lint
make tidy           # go mod tidy
make clean          # Remove build artifacts
make install        # Install to /usr/local/bin
make dev ARGS="..." # Run with go run
```

## Running Tests

```bash
# All tests
make test

# Verbose
go test -v ./...

# Specific package
go test -v ./internal/project/...
go test -v ./internal/share/...
go test -v ./internal/ziputil/...

# With coverage
make test-coverage
# Opens coverage.html in browser

# Race detector
go test -race ./...
```

## Code Style

### Formatting

```bash
# Auto-format
make fmt
# or
go fmt ./...
```

### Linting

```bash
# Run linter
make lint
# or
golangci-lint run
```

### Key Style Rules

- Use `gofmt` formatting (enforced by `make fmt`)
- Follow Go idioms and best practices
- Error handling: always check errors
- Use meaningful variable names
- Keep functions small and focused
- Write tests for new functionality

## Adding a New Command

### 1. Add Command to CLI App

Edit `internal/cli/app.go`:

```go
func (a *App) Run(args []string) error {
    // ... existing switch ...
    case "mycommand":
        return a.runMyCommand(args[2:])
    // ...
}

func (a *App) runMyCommand(args []string) error {
    fs := flag.NewFlagSet("mycommand", flag.ExitOnError)
    // Define flags
    fs.Parse(args)
    
    // Implementation
    return nil
}
```

### 2. Add Help Text

Update `printUsage()` in `internal/cli/app.go`.

### 3. Add Tests

Create `internal/cli/app_test.go`:
```go
func TestApp_RunMyCommand(t *testing.T) {
    app := NewApp()
    err := app.runMyCommand([]string{"-flag", "value"})
    assert.NoError(t, err)
}
```

### 4. Update Documentation

- `wiki/Commands-Reference.md`
- `wiki/Quick-Start.md` (if user-facing)
- `README.md`

## Adding a New Package

### 1. Create Package Directory

```bash
mkdir -p internal/newpackage
```

### 2. Create Package Files

`internal/newpackage/manager.go`:
```go
package newpackage

type Manager struct{}

func NewManager() *Manager {
    return &Manager{}
}

func (m *Manager) DoSomething() error {
    return nil
}
```

`internal/newpackage/manager_test.go`:
```go
package newpackage

import "testing"

func TestManager_DoSomething(t *testing.T) {
    m := NewManager()
    err := m.DoSomething()
    if err != nil {
        t.Fatal(err)
    }
}
```

### 3. Import in CLI

```go
import "GoShare/internal/newpackage"

func NewApp() *App {
    return &App{
        // ...
        newPackage: newpackage.NewManager(),
    }
}
```

## Architecture Guidelines

### Package Responsibilities

| Package | Responsibility |
|---------|----------------|
| `cli` | Command parsing, user interaction |
| `project` | Project CRUD, file tree, HTML generation |
| `publish` | Cloudflare tunnel, HTTP serving |
| `share` | Share codes, ZIP creation, download |
| `ziputil` | ZIP compression, extraction, security |

### Dependency Rules

```
cli → project, publish, share, ziputil
project → ziputil (for copy)
share → ziputil
publish → (stdlib only)
ziputil → (stdlib only)
```

**No circular dependencies!**

### Error Handling

```go
// Good: wrap errors with context
if err := os.MkdirAll(path, 0755); err != nil {
    return fmt.Errorf("failed to create directory %s: %w", path, err)
}

// Good: sentinel errors for expected cases
var ErrProjectNotFound = errors.New("project not found")

// Avoid: ignoring errors
os.MkdirAll(path, 0755) // BAD
```

### Logging

Use `fmt` for user output, `log` for debug:
```go
// User-facing
color.Green("Success!")
color.Red("Error: %v", err)

// Debug (optional)
log.Printf("DEBUG: processing %s", path)
```

## Testing Guidelines

### Unit Tests

- Test one thing per test function
- Use table-driven tests for multiple cases
- Mock external dependencies
- Test error paths

```go
func TestManager_Create(t *testing.T) {
    tests := []struct {
        name        string
        projectName string
        summary     string
        wantErr     bool
    }{
        {"valid", "test", "summary", false},
        {"empty name", "", "summary", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test...
        })
    }
}
```

### Integration Tests

- Test full command flows
- Use temporary directories
- Clean up in `defer`

```go
func TestIntegration_CreateAndShare(t *testing.T) {
    tmpDir := t.TempDir()
    oldHome := os.Getenv("HOME")
    os.Setenv("HOME", tmpDir)
    defer os.Setenv("HOME", oldHome)
    
    // Test full flow...
}
```

## Debugging

### Build with Debug Info

```bash
go build -gcflags="all=-N -l" -o bin/goshare-debug ./cmd/goshare
```

### Delve Debugger

```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./cmd/goshare -- create -project test
```

### Print Debug Info

```bash
# Verbose output
go run ./cmd/goshare create -project test 2>&1 | head -20
```

## Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof ./internal/project/...
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof ./internal/project/...
go tool pprof mem.prof
```

## Release Process

### 1. Update Version

```bash
# Edit Makefile
VERSION=1.1.0
```

### 2. Update Changelog

```markdown
## [1.1.0] - 2024-01-15
### Added
- New feature X
### Fixed
- Bug Y
```

### 3. Build Release Binaries

```bash
make build-all
```

### 4. Create Git Tag

```bash
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0
```

### 5. GitHub Release

- Go to GitHub Releases
- Create new release from tag
- Upload binaries from `bin/`
- Generate release notes

## Contributing Guidelines

### Pull Request Checklist

- [ ] Tests pass (`make test`)
- [ ] Code formatted (`make fmt`)
- [ ] Linter passes (`make lint`)
- [ ] Vet passes (`make vet`)
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] No breaking changes (or documented)

### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Examples:
```
feat(cli): add new list command
fix(ziputil): handle path traversal in extract
docs(wiki): update sharing guide
test(project): add create duplicate test
```

## Useful Commands

```bash
# Quick test cycle
make dev ARGS="create -project test"

# Watch for changes (requires entr)
ls **/*.go | entr -c make test

# Generate mocks (if using mockery)
go install github.com/vektra/mockery/v2@latest
mockery --name=Manager --dir=internal/project --output=internal/project/mocks

# Update dependencies
go get -u ./...
go mod tidy

# Check for vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## IDE Setup

### VS Code

Install Go extension, then configure:
```json
// .vscode/settings.json
{
    "go.formatTool": "gofmt",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v"],
    "go.buildFlags": ["-ldflags=-X main.version=dev"]
}
```

### GoLand

- Enable "Run gofmt on save"
- Configure golangci-lint as external tool
- Set up run configuration for `cmd/goshare/main.go`

## Performance Tips

1. **Avoid repeated filesystem walks** - Cache file trees
2. **Use buffered I/O** for large file operations
3. **Stream ZIP extraction** - Don't load entire archive in memory
4. **Reuse HTTP clients** - For cloudflared communication

## Security Checklist

- [ ] Validate all file paths (path traversal)
- [ ] Sanitize user input in templates
- [ ] Use crypto/rand for random generation
- [ ] Set secure file permissions (0644, 0755)
- [ ] No hardcoded secrets
- [ ] Dependencies up to date (`govulncheck`)

---

## Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Project Layout](https://github.com/golang-standards/project-layout)
- [Cobra CLI](https://cobra.dev/) (for reference, we use flag package)