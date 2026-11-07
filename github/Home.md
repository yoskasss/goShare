# GoShare Wiki

Welcome to the GoShare documentation wiki. GoShare is a lightweight Go command-line tool for creating, managing, previewing, and sharing project directories.

## Quick Links

- [Installation](Installation.md)
- [Quick Start](Quick-Start.md)
- [Commands Reference](Commands-Reference.md)
- [Project Structure](Project-Structure.md)
- [Sharing Projects](Sharing-Projects.md)
- [Publishing with Cloudflare](Publishing-with-Cloudflare.md)
- [Configuration](Configuration.md)
- [Development](Development.md)
- [Troubleshooting](Troubleshooting.md)
- [FAQ](FAQ.md)

## Overview

GoShare helps you:
- 📁 **Create** project directories with metadata
- 📋 **Generate** GitHub-style HTML project pages
- 📦 **Copy** files/directories into projects
- 🔗 **Share** projects via short codes
- ☁️ **Publish** projects publicly via Cloudflare tunnels
- 🗑️ **Manage** project lifecycle (list, delete)

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Create    │────▶│   Project   │────▶│  index.html │
│  Project    │     │  Directory  │     │  (GitHub    │
└─────────────┘     └─────────────┘     │   style)    │
       │                   │            └─────────────┘
       │                   │
       ▼                   ▼
┌─────────────┐     ┌─────────────┐
│   Copy      │     │   Share     │
│  Files      │     │  (ZIP +     │
└─────────────┘     │   Code)     │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Download   │
                    │  & Extract  │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Publish    │
                    │ (Cloudflare)│
                    └─────────────┘
```

## Data Storage

| Data | Location | Description |
|------|----------|-------------|
| Projects | `~/GoShareProjects/` | Project directories with metadata |
| Shares | `~/GoShareShares/` | ZIP archives + share codes (7-day expiry) |

## Requirements

- **Go 1.25+** (for building from source)
- **cloudflared** (for publishing - [install guide](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation))

## License

MIT License - see [LICENSE](../LICENSE) for details.