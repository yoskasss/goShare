# Configuration

GoShare configuration options and customization.

## Configuration Files

GoShare doesn't use a central config file. Configuration is done via:
1. **Command-line flags** (primary)
2. **Environment variables** (limited)
3. **Source code modification** (advanced)

## Command-Line Flags

### Global Flags (all commands)

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `-project` | `-p` | Project name | `GoShare` |
| `-help` | `-h` | Show command help | - |

### Command-Specific Flags

#### `create`
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `-project` | `-p` | Project name | `GoShare` |
| `-summary` | - | Project summary | (prompts) |

#### `copy` / `cp`
| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| `-project` | `-p` | Target project | Yes |
| `-cp` | `-c` | Source path | Yes |

#### `public`
| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| `-project` | `-p` | Project to publish | Yes |

#### `share`
| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| `-project` | `-p` | Project to share | Yes |

#### `download`
| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| `-share` | `-s` | Share code | Yes |
| `-setup` | - | Custom destination | No |

#### `delete`
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `-project` | `-p` | Project to delete | Yes |
| `-force` | `-f` | Skip confirmation | `false` |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `HOME` | User home directory (data paths) | OS default |
| `BROWSER` | Browser to open index.html | System default |

### Setting Environment Variables

```bash
# Linux/macOS (bash/zsh)
export BROWSER=firefox
export HOME=/custom/home

# Windows (PowerShell)
$env:BROWSER = "chrome"
$env:HOME = "C:\Custom\Home"

# Persistent (add to ~/.bashrc, ~/.zshrc, or PowerShell profile)
echo 'export BROWSER=firefox' >> ~/.bashrc
```

## Data Directory Locations

### Default Locations

| Data | Location |
|------|----------|
| Projects | `~/GoShareProjects/` |
| Shares | `~/GoShareShares/` |

### Custom Locations (Source Modification)

Edit `internal/project/manager.go`:
```go
func NewManager() *Manager {
    homeDir, _ := os.UserHomeDir()
    return &Manager{
        basePath: filepath.Join(homeDir, "CustomProjects"),  // Change here
    }
}
```

Edit `internal/share/manager.go`:
```go
func NewManager() *Manager {
    homeDir, _ := os.UserHomeDir()
    return &Manager{
        shareDir: filepath.Join(homeDir, "CustomShares"),  // Change here
    }
}
```

Then rebuild: `make build`

## Project Template Customization

### Modify index.html Template

Edit `internal/project/manager.go` → `generateIndexHTML()` function:

```go
tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Name}} - My Custom Template</title>
    <style>
        /* Your custom CSS */
        body { font-family: 'Your Font', sans-serif; }
        .header { background: #your-color; }
    </style>
</head>
<body>
    <!-- Your custom HTML -->
</body>
</html>`
```

### Template Variables

| Variable | Type | Description |
|----------|------|-------------|
| `{{.Name}}` | string | Project name |
| `{{.Summary}}` | string | Project summary |
| `{{.CreatedAt}}` | time.Time | Creation timestamp |
| `{{.UpdatedAt}}` | time.Time | Last update timestamp |
| `{{.FileTree}}` | string | Generated file tree |
| `{{.Path}}` | string | Project path |

### Custom Template Example

```go
tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Name}}</title>
    <style>
        body { max-width: 800px; margin: 0 auto; padding: 2rem; font-family: system-ui; }
        .badge { background: #6366f1; color: white; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; }
        pre { background: #1e1e1e; color: #d4d4d4; padding: 1rem; border-radius: 8px; overflow-x: auto; }
    </style>
</head>
<body>
    <header>
        <h1>{{.Name}} <span class="badge">GoShare</span></h1>
        <p>{{.Summary}}</p>
        <small>Created: {{.CreatedAt.Format "Jan 2, 2006"}} | Updated: {{.UpdatedAt.Format "Jan 2, 2006"}}</small>
    </header>
    <section>
        <h2>Files</h2>
        <pre>{{.FileTree}}</pre>
    </section>
</body>
</html>`
```

## Share Configuration

### Expiration Time

Edit `internal/share/manager.go`:
```go
// Change 7 days to desired duration
ExpiresAt: time.Now().Add(30 * 24 * time.Hour),  // 30 days
```

### Share Code Length

```go
// Change 6 bytes (8 chars) to desired length
bytes := make([]byte, 8)  // 11 chars
// or
bytes := make([]byte, 4)  // 6 chars
```

### Share Directory

```go
func NewManager() *Manager {
    return &Manager{
        shareDir: "/custom/path/shares",  // Custom location
    }
}
```

## ZIP Configuration

### Compression Level

Edit `internal/ziputil/manager.go`:
```go
zipWriter := zip.NewWriter(zipFile)
// Add: zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
//     return flate.NewWriter(out, flate.BestCompression)
// })
```

### Exclude Patterns

Modify `CreateZip` function:
```go
// Add more exclusions
excludedPrefixes := []string{".goshare.json", ".git", ".env", "node_modules", "vendor"}
for _, prefix := range excludedPrefixes {
    if strings.Contains(relPath, prefix) {
        return nil
    }
}
```

## Shell Completion Configuration

### Generate Completion Scripts

```bash
# Bash
goshare completion bash > /etc/bash_completion.d/goshare

# Zsh
goshare completion zsh > /usr/local/share/zsh/site-functions/_goshare

# Fish
goshare completion fish > ~/.config/fish/completions/goshare.fish

# PowerShell
goshare completion powershell > goshare.ps1
# Then source in profile
```

### Custom Completion (Advanced)

The completion is auto-generated from Cobra flags. To customize, modify `internal/cli/app.go` flag definitions.

## Build-Time Configuration

### Version Information

Set via Makefile LDFLAGS:
```makefile
VERSION=1.0.0
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date -u +%Y-%m-%d)"
```

Add to `cmd/goshare/main.go`:
```go
var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    // Add version command
    if len(os.Args) > 1 && os.Args[1] == "version" {
        fmt.Printf("GoShare %s (commit: %s, built: %s)\n", version, commit, date)
        return
    }
    // ...
}
```

### Build Tags

```bash
# Build with custom tags
go build -tags "custom" ./cmd/goshare
```

## Runtime Configuration (Planned)

Future versions may support:
- `~/.config/goshare/config.yaml`
- Project-specific `.goshare.yaml` in project root
- Global defaults for flags

## Example: Complete Custom Setup

```bash
#!/bin/bash
# custom-setup.sh

# 1. Set custom directories
export GOSHARE_PROJECTS="$HOME/Dev/GoShareProjects"
export GOSHARE_SHARES="$HOME/Dev/GoShareShares"

# 2. Create wrapper script
cat > ~/bin/goshare-custom << 'EOF'
#!/bin/bash
export HOME="$HOME/Dev"
exec /usr/local/bin/goshare "$@"
EOF
chmod +x ~/bin/goshare-custom

# 3. Add aliases
cat >> ~/.bashrc << 'EOF'
alias gs="goshare-custom"
alias gsc="goshare-custom create"
alias gscp="goshare-custom copy"
alias gss="goshare-custom share"
alias gsd="goshare-custom download"
alias gsl="goshare-custom list"
EOF

# 4. Reload
source ~/.bashrc
```

## Configuration Reference Summary

| Setting | Method | Location |
|---------|--------|----------|
| Project name | Flag | `-project` / `-p` |
| Project summary | Flag | `-summary` (create) |
| Source path | Flag | `-cp` / `-c` (copy) |
| Share code | Flag | `-share` / `-s` (download) |
| Custom destination | Flag | `-setup` (download) |
| Force delete | Flag | `-force` / `-f` (delete) |
| Data directories | Source | `manager.go` constructors |
| HTML template | Source | `manager.go` generateIndexHTML() |
| Share expiry | Source | `manager.go` GenerateShareCode() |
| Share code length | Source | `manager.go` generateCode() |
| Exclude patterns | Source | `ziputil/manager.go` CreateZip() |
| Version info | Build | Makefile LDFLAGS |