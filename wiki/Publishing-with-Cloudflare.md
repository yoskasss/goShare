# Publishing with Cloudflare

Guide to publishing projects publicly using Cloudflare tunnels.

## Overview

The `public` command:
1. Starts a local HTTP server on port 8080
2. Creates a Cloudflare Quick Tunnel
3. Provides a public HTTPS URL
4. Serves your entire project directory

## Prerequisites

### Install cloudflared

```bash
# Linux (Debian/Ubuntu)
curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared.deb

# macOS
brew install cloudflared

# Windows (Scoop)
scoop install cloudflared

# Verify
cloudflared version
```

### Verify Installation

```bash
cloudflared version
# cloudflared version 2024.x.x (built 2024-xx-xx-xx)
```

## Basic Usage

```bash
# Publish a project
goshare public -project my-project
```

### Output

```
Starting cloudflared tunnel for project 'my-project'...
Press Ctrl+C to stop the tunnel
Tunnel started! Share the URL shown above.
Serving files from: /home/user/GoShareProjects/my-project
```

Cloudflare will display:
```
https://random-adjective-noun.trycloudflare.com
```

Share this URL with anyone - they can view your project in their browser.

### Stop the Tunnel

Press `Ctrl+C` to stop both the local server and Cloudflare tunnel.

## What Gets Served

The entire project directory is served statically:

```
https://xxx.trycloudflare.com/           → index.html (project page)
https://xxx.trycloudflare.com/src/       → Directory listing
https://xxx.trycloudflare.com/src/main.go → Raw file
https://xxx.trycloudflare.com/docs/      → Directory listing
https://xxx.trycloudflare.com/image.png  → Image file
```

## Advanced Usage

### Custom Port (if 8080 is busy)

Currently hardcoded to 8080. To change, modify `internal/publish/publisher.go`.

### Run in Background

```bash
# Start in background
nohup goshare public -project my-project > tunnel.log 2>&1 &

# Get the URL from logs
grep "trycloudflare.com" tunnel.log

# Stop later
pkill -f "goshare public"
```

### With Custom Domain (Cloudflare Tunnel)

For production use, create a named tunnel:

```bash
# 1. Login to Cloudflare
cloudflared tunnel login

# 2. Create tunnel
cloudflared tunnel create my-project-tunnel

# 3. Configure DNS
cloudflared tunnel route dns my-project-tunnel project.example.com

# 4. Run tunnel (in background)
cloudflared tunnel run my-project-tunnel

# 5. Start local server only
cd ~/GoShareProjects/my-project
python3 -m http.server 8080
```

Then access at `https://project.example.com`

## Troubleshooting

### "cloudflared not found"

```bash
# Install cloudflared
# See Prerequisites section above

# Verify it's in PATH
which cloudflared
cloudflared version
```

### "Port 8080 already in use"

```bash
# Find what's using port 8080
lsof -i :8080
# or
netstat -tulpn | grep 8080

# Kill the process
kill -9 <PID>

# Or use different port (modify source)
```

### Tunnel URL not appearing

```bash
# Check cloudflared logs
# Run with more verbosity
cloudflared tunnel --url http://localhost:8080 --loglevel debug
```

### "Connection refused" / "502 Bad Gateway"

```bash
# Check local server is running
curl http://localhost:8080

# Should return HTML content
# If not, check project has index.html
ls -la ~/GoShareProjects/my-project/index.html
```

### Firewall/Network Issues

```bash
# Test local access first
curl http://localhost:8080

# Check if cloudflared can reach Cloudflare
cloudflared tunnel --url http://localhost:8080 --hello-world
```

## Security Considerations

### What's Exposed

⚠️ **Everything in the project directory is publicly accessible**

- Source code
- Configuration files
- Documentation
- Images, assets
- Any file you copied into the project

### Best Practices

1. **Review before publishing** - Check for secrets, API keys, passwords
2. **Use `.gitignore` patterns** - Don't copy sensitive files
3. **Limit exposure time** - Stop tunnel when done
4. **Monitor access** - Check Cloudflare dashboard for analytics

### What's NOT Exposed

- `.goshare.json` (metadata file)
- Hidden files (starting with `.`)
- Your local filesystem outside the project directory

## Cloudflare Dashboard

View tunnel metrics at: https://dash.cloudflare.com

1. Go to **Zero Trust** → **Networks** → **Tunnels**
2. Find your quick tunnel (random name)
3. View:
   - Request count
   - Bandwidth
   - Errors
   - Visitor locations

## Automation

### Publish in CI/CD

```yaml
# .github/workflows/publish.yml
name: Publish Preview
on:
  pull_request:
    types: [opened, synchronize]
jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: make build
      - run: |
          ./bin/goshare create -project preview-${{ github.event.number }} \
            -summary "PR #${{ github.event.number }} preview"
      - run: ./bin/goshare cp -project preview-${{ github.event.number }} -cp .
      - name: Start tunnel
        run: |
          timeout 300 ./bin/goshare public -project preview-${{ github.event.number }} &
          sleep 10
          # Extract URL from logs (simplified)
          echo "Preview URL available in logs"
```

### Scheduled Publishing

```bash
#!/bin/bash
# publish-daily.sh

PROJECT="daily-build"
goshare create -project "$PROJECT" -summary "Daily build $(date)" 2>/dev/null || true
goshare cp -project "$PROJECT" -cp ./dist

# Run for 1 hour
timeout 3600 goshare public -project "$PROJECT"
```

## Alternatives

If Cloudflare doesn't work for you:

| Tool | Command | Notes |
|------|---------|-------|
| ngrok | `ngrok http 8080` | Requires account, free tier limited |
| localtunnel | `npx localtunnel --port 8080` | No account needed, random subdomain |
| serveo | `ssh -R 80:localhost:8080 serveo.net` | SSH-based, no install |
| VS Code Port Forwarding | Built-in | If using VS Code Remote |

## FAQ

**Q: Does the tunnel persist after I close the terminal?**
A: No, it stops when you press Ctrl+C or close the terminal.

**Q: Can I use a custom domain?**
A: Yes, with a named Cloudflare Tunnel (see Advanced Usage).

**Q: Is there a time limit on quick tunnels?**
A: Quick tunnels run until stopped. Named tunnels can run indefinitely.

**Q: Can multiple people access simultaneously?**
A: Yes, Cloudflare handles concurrent connections.

**Q: Does it work behind corporate firewalls?**
A: Usually yes, Cloudflare uses standard HTTPS (port 443).

**Q: Can I password-protect the tunnel?**
A: Not with quick tunnels. Use Cloudflare Access with named tunnels for authentication.

**Q: What's the bandwidth limit?**
A: Free tier has generous limits. Check Cloudflare pricing for details.