# Troubleshooting

Common issues and solutions for GoShare.

## Installation Issues

### "command not found: goshare"

**Cause:** Binary not in PATH.

**Solution:**
```bash
# Check where binary is
ls -la ~/GoShare/bin/goshare
# or
ls -la /usr/local/bin/goshare

# Add to PATH
export PATH=$PATH:~/GoShare/bin
# or
sudo ln -s ~/GoShare/bin/goshare /usr/local/bin/goshare
```

### "go: cannot find main module" during build

**Cause:** Not in project root or go.mod missing.

**Solution:**
```bash
cd /path/to/GoShare
ls go.mod  # Should exist
make build
```

### "cloudflared: command not found" when publishing

**Cause:** cloudflared not installed or not in PATH.

**Solution:**
```bash
# Install cloudflared
# See Installation wiki page

# Verify
which cloudflared
cloudflared version
```

## Project Issues

### "project 'xxx' already exists"

**Cause:** Project with that name already exists.

**Solution:**
```bash
# List existing projects
goshare list

# Delete old project
goshare delete -project xxx -force

# Or use different name
goshare create -project xxx-new
```

### "project 'xxx' not found"

**Cause:** Project doesn't exist or wrong name.

**Solution:**
```bash
# Check exact name (case-sensitive)
goshare list

# Check filesystem
ls ~/GoShareProjects/
```

### "index.html not found in project"

**Cause:** Project corrupted or manually deleted index.html.

**Solution:**
```bash
# Recreate project (copies files over)
goshare create -project xxx -summary "Restored"

# Or manually regenerate (copy files first)
cp -r ~/GoShareProjects/xxx ~/GoShareProjects/xxx-backup
goshare delete -project xxx -force
goshare create -project xxx -summary "Restored"
goshare cp -project xxx -cp ~/GoShareProjects/xxx-backup
```

### Permission denied when creating project

**Cause:** No write access to home directory.

**Solution:**
```bash
# Check permissions
ls -ld ~/

# Fix if needed
chmod 755 ~/

# Or use custom location (modify source)
```

## Copy Issues

### "source not found"

**Cause:** Path doesn't exist or typo.

**Solution:**
```bash
# Use absolute path
goshare cp -project xxx -cp /full/path/to/source

# Check source exists
ls -la /path/to/source

# Tab completion
goshare cp -project xxx -cp ./src<TAB>
```

### "failed to copy directory"

**Cause:** Permission issue, disk full, or special files.

**Solution:**
```bash
# Check disk space
df -h ~/

# Check permissions
ls -ld /source/path

# Try copying manually first
cp -r /source/path ~/GoShareProjects/xxx/
```

### Copied files not showing in index.html

**Cause:** index.html not regenerated.

**Solution:**
```bash
# Copy command should regenerate automatically
# If not, recreate project page
goshare create -project xxx -summary "Existing project"
# This preserves files but regenerates index.html
```

## Share Issues

### "share code 'xxx' not found or expired"

**Cause:** Code wrong, expired (7 days), or files deleted.

**Solution:**
```bash
# Check if share files exist
ls -la ~/GoShareShares/xxx.*

# Check expiration
cat ~/GoShareShares/xxx.json | jq '.expires_at'
# Compare with: date -u +"%Y-%m-%dT%H:%M:%SZ"

# Regenerate share
goshare share -project my-project
```

### "archive not found for share code"

**Cause:** ZIP file deleted but JSON remains.

**Solution:**
```bash
# Clean up orphaned JSON
rm ~/GoShareShares/xxx.json

# Regenerate
goshare share -project my-project
```

### "invalid share info"

**Cause:** JSON corrupted.

**Solution:**
```bash
# Remove both files
rm ~/GoShareShares/xxx.json ~/GoShareShares/xxx.zip

# Regenerate
goshare share -project my-project
```

### Download extracts to wrong location

**Cause:** Default behavior uses project name.

**Solution:**
```bash
# Use -setup flag
goshare download -share xxx -setup /exact/path

# Or rename after download
mv ~/GoShareProjects/old-name ~/GoShareProjects/new-name
```

### "downloaded file not found" after download

**Cause:** Project name mismatch or extraction issue.

**Solution:**
```bash
# Check what was extracted
ls -la ~/GoShareProjects/

# Check share metadata
cat ~/GoShareShares/xxx.json | jq '.project_name'

# Download with explicit destination
goshare download -share xxx -setup ~/my-custom-path
```

## Publish Issues

### "cloudflared not found in PATH"

**Cause:** cloudflared not installed.

**Solution:**
```bash
# Install cloudflared
# See Publishing wiki page

# Verify
cloudflared version
```

### "port 8080 already in use"

**Cause:** Another process using port 8080.

**Solution:**
```bash
# Find process
lsof -i :8080
# or
ss -tulpn | grep 8080

# Kill it
kill -9 <PID>

# Or wait for it to release
```

### Tunnel URL not displayed

**Cause:** cloudflared output buffering or error.

**Solution:**
```bash
# Run cloudflared directly to debug
cd ~/GoShareProjects/my-project
cloudflared tunnel --url http://localhost:8080

# Check for errors in output
```

### "connection refused" when accessing tunnel URL

**Cause:** Local server not running or wrong port.

**Solution:**
```bash
# Test local server
curl http://localhost:8080

# Should return HTML
# If not, check project has index.html
ls ~/GoShareProjects/my-project/index.html
```

### Tunnel works but shows 404 for files

**Cause:** File paths or server configuration.

**Solution:**
```bash
# Check file exists
ls ~/GoShareProjects/my-project/src/main.go

# Test direct access
curl http://localhost:8080/src/main.go

# Should return file content
```

## General Issues

### High memory usage

**Cause:** Large projects or many shares.

**Solution:**
```bash
# Clean old shares
find ~/GoShareShares -mtime +7 -delete

# Check project sizes
du -sh ~/GoShareProjects/*

# Remove unused projects
goshare delete -project old-project -force
```

### Slow performance

**Cause:** Large file trees, slow disk.

**Solution:**
```bash
# Exclude large directories from copy
# Use .gitignore patterns (manual)

# Check disk speed
dd if=/dev/zero of=~/test bs=1M count=1024 conv=fdatasync
```

### "too many open files"

**Cause:** System limit reached.

**Solution:**
```bash
# Check limit
ulimit -n

# Increase temporarily
ulimit -n 4096

# Permanent: edit /etc/security/limits.conf
```

## Debug Mode

### Enable Verbose Output

```bash
# GoShare doesn't have a debug flag yet
# But you can run with go run for more output
go run ./cmd/goshare create -project test
```

### Check Logs

```bash
# No persistent logs by default
# For cloudflared, run manually:
cloudflared tunnel --url http://localhost:8080 --loglevel debug
```

### Inspect Data Files

```bash
# Project metadata
cat ~/GoShareProjects/my-project/.goshare.json | jq .

# Share metadata
cat ~/GoShareShares/xxx.json | jq .

# ZIP contents
unzip -l ~/GoShareShares/xxx.zip
```

## Getting Help

### Check Version

```bash
goshare version
# GoShare v1.0.0
```

### Run Self-Test

```bash
# Run all tests
make test

# Should pass all tests
```

### Report Issues

When reporting issues, include:

1. **OS and version:** `uname -a` or Windows version
2. **GoShare version:** `goshare version`
3. **Go version:** `go version`
4. **Command that failed:** Full command with flags
5. **Error message:** Complete output
6. **Steps to reproduce:** Minimal example

### Common Diagnostic Commands

```bash
# Environment
echo "HOME=$HOME"
echo "PATH=$PATH"
go version
goshare version
cloudflared version

# Data directories
ls -la ~/GoShareProjects/
ls -la ~/GoShareShares/

# Permissions
ls -ld ~/GoShareProjects ~/GoShareShares

# Disk space
df -h ~/
```

## Error Code Reference

| Error | Code | Solution |
|-------|------|----------|
| Project exists | 1 | Delete or rename |
| Project not found | 1 | Check name, list projects |
| Source not found | 1 | Check path, use absolute |
| Share expired | 1 | Regenerate share |
| Cloudflared missing | 1 | Install cloudflared |
| Port in use | 1 | Kill process on 8080 |
| Permission denied | 1 | Check directory permissions |
| Disk full | 1 | Clean up, check df -h |

## Reset Everything

```bash
# WARNING: Deletes all projects and shares!
rm -rf ~/GoShareProjects ~/GoShareShares

# Then recreate
goshare create -project new-project
```