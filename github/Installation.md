# Installation

## Pre-built Binaries (Recommended)

Download the latest release for your platform:

| Platform | Architecture | Download |
|----------|--------------|----------|
| Linux | x86_64 | `goshare-linux-amd64` |
| macOS | x86_64 | `goshare-darwin-amd64` |
| macOS | ARM64 (M1/M2/M3) | `goshare-darwin-arm64` |
| Windows | x86_64 | `goshare-windows-amd64.exe` |

### Linux/macOS

```bash
# Download (replace with actual release URL)
curl -L -o goshare https://github.com/yourusername/GoShare/releases/latest/download/goshare-linux-amd64

# Make executable
chmod +x goshare

# Move to PATH
sudo mv goshare /usr/local/bin/

# Verify
goshare version
```

### Windows (PowerShell)

```powershell
# Download
Invoke-WebRequest -Uri "https://github.com/yourusername/GoShare/releases/latest/download/goshare-windows-amd64.exe" -OutFile "goshare.exe"

# Move to a directory in your PATH
Move-Item goshare.exe "C:\Program Files\GoShare\goshare.exe"

# Add to PATH (if not already)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\GoShare\", "User")

# Verify
goshare version
```

## From Source

### Prerequisites

- Go 1.25 or newer
- Git

### Build Steps

```bash
# Clone repository
git clone https://github.com/yourusername/GoShare.git
cd GoShare

# Build
make build

# Binary is at ./bin/goshare
./bin/goshare version

# Optional: Install to PATH
make install
```

### Cross-Platform Builds

```bash
# Build for all platforms
make build-all

# Output in ./bin/
ls -la bin/
# goshare-linux-amd64
# goshare-darwin-amd64
# goshare-darwin-arm64
# goshare-windows-amd64.exe
```

## Install cloudflared (for Publishing)

Publishing requires `cloudflared` to be installed and in your PATH.

### Linux

```bash
# Debian/Ubuntu
curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared.deb

# RHEL/CentOS/Fedora
curl -L --output cloudflared.rpm https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-x86_64.rpm
sudo rpm -i cloudflared.rpm

# Arch Linux
sudo pacman -S cloudflared
```

### macOS

```bash
# Homebrew
brew install cloudflared

# Or direct download
curl -L --output cloudflared.pkg https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-darwin-amd64.pkg
sudo installer -pkg cloudflared.pkg -target /
```

### Windows

```powershell
# Scoop
scoop install cloudflared

# Chocolatey
choco install cloudflared

# Or download .msi from GitHub releases
```

### Verify Installation

```bash
cloudflared version
# Should output something like: cloudflared version 2024.x.x (built 2024-xx-xx-xx)
```

## Docker (Alternative)

```dockerfile
# Dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o goshare ./cmd/goshare

FROM alpine:latest
RUN apk --no-cache add ca-certificates cloudflared
COPY --from=builder /app/goshare /usr/local/bin/
ENTRYPOINT ["goshare"]
```

```bash
# Build
docker build -t goshare .

# Run
docker run --rm -v ~/GoShareProjects:/root/GoShareProjects goshare create -project myproject
```

## Shell Completion

### Bash

```bash
# Generate completion script
goshare completion bash > /etc/bash_completion.d/goshare

# Or for current user
goshare completion bash > ~/.bash_completion.d/goshare
source ~/.bash_completion.d/goshare
```

### Zsh

```bash
goshare completion zsh > "${fpath[1]}/_goshare"
# Restart shell or run: compinit
```

### Fish

```bash
goshare completion fish > ~/.config/fish/completions/goshare.fish
```

### PowerShell

```powershell
goshare completion powershell | Out-String | Invoke-Expression
# Add to profile for persistence
goshare completion powershell >> $PROFILE
```

## Verify Installation

```bash
# Check version
goshare version
# GoShare v1.0.0

# Check help
goshare help

# Check cloudflared (for publishing)
cloudflared version
```

## Uninstall

```bash
# Remove binary
sudo rm /usr/local/bin/goshare

# Remove data directories (optional)
rm -rf ~/GoShareProjects ~/GoShareShares

# Remove shell completions
rm /etc/bash_completion.d/goshare  # or user-specific path
```