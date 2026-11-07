# Sharing Projects

Complete guide to sharing projects with GoShare.

## Overview

GoShare's sharing system works by:
1. **Creating a ZIP archive** of your project
2. **Generating a short code** (8 characters)
3. **Storing both** locally with 7-day expiry
4. **Recipient downloads** using the code

## Creating a Share

```bash
goshare share -project my-project
```

Output:
```
Creating archive for sharing...
Share code: aB3x9K2m
Download with: goshare download -share aB3x9K2m
```

### What Gets Shared

| Included | Excluded |
|----------|----------|
| All project files | `.goshare.json` (metadata) |
| Directory structure | Hidden files (`.git`, `.env`, etc.) |
| `index.html` | Build artifacts (if hidden) |
| Source code, docs, images | - |

### Share Code Properties

- **Format:** 8-character URL-safe base64
- **Example:** `aB3x9K2m`, `Xy7zP2qR`, `h64bwBkf`
- **Case-sensitive:** Yes
- **Lifetime:** 7 days
- **One-time use:** No, can be used multiple times

## Sharing the Code

### Secure Sharing Methods

| Method | Security | Use Case |
|--------|----------|----------|
| Password manager (Bitwarden, 1Password) | 🔒 High | Team sharing |
| Encrypted chat (Signal, Wire) | 🔒 High | Private sharing |
| Email (PGP encrypted) | 🔒 Medium | External partners |
| Plain text chat/email | ⚠️ Low | Temporary/internal only |

### Sharing Template

```
Hey [name],

I've shared a GoShare project with you:
- **Project:** my-awesome-project
- **Share code:** `aB3x9K2m`
- **Expires:** 7 days from now

To download:
```bash
goshare download -share aB3x9K2m
```

Let me know if you have issues!
```

## Downloading a Shared Project

### Basic Download

```bash
goshare download -share aB3x9K2m
```

Downloads to `~/GoShareProjects/<original-project-name>/`

### Custom Destination

```bash
# Specific directory
goshare download -share aB3x9K2m -setup /tmp/my-project

# Current directory
goshare download -share aB3x9K2m -setup .

# Relative path
goshare download -share aB3x9K2m -setup ./restored-project
```

### What Happens During Download

1. Validates share code exists
2. Checks expiration (7 days)
3. Extracts ZIP to destination
4. Creates project directory with original name
5. Shows success message with path

## Advanced Sharing Workflows

### Share Multiple Projects

```bash
# Create shares for multiple projects
goshare share -project frontend
goshare share -project backend
goshare share -project docs

# Share all codes at once
echo "Frontend: aB3x9K2m"
echo "Backend: Xy7zP2qR"
echo "Docs: h64bwBkf"
```

### Automated Sharing (CI/CD)

```bash
#!/bin/bash
# share-on-release.sh

PROJECT="my-app"
VERSION=$(git describe --tags)

# Create project if needed
goshare create -project "$PROJECT" -summary "Release $VERSION" 2>/dev/null || true

# Copy build artifacts
goshare cp -project "$PROJECT" -cp ./dist

# Generate share code
CODE=$(goshare share -project "$PROJECT" 2>&1 | grep "Share code:" | awk '{print $3}')

echo "Release $VERSION share code: $CODE"
echo "Download: goshare download -share $CODE"

# Post to Slack, Discord, etc.
curl -X POST -H 'Content-type: application/json' \
  --data "{\"text\":\"New release $VERSION: $CODE\"}" \
  $WEBHOOK_URL
```

### Share with Expiration Notice

```bash
# Check share status before sending
goshare share-info aB3x9K2m  # Planned feature

# Or manually check
cat ~/GoShareShares/aB3x9K2m.json | jq '.expires_at'
```

## Managing Shares

### List Active Shares

```bash
# Not directly supported yet, but you can check:
ls -la ~/GoShareShares/*.json
```

### Clean Up Expired Shares

```bash
# Automatic: runs on share/download commands
# Manual (planned):
goshare share-cleanup

# Or manually:
find ~/GoShareShares -name "*.json" -mtime +7 -delete
find ~/GoShareShares -name "*.zip" -mtime +7 -delete
```

### Delete a Specific Share

```bash
# Remove share code and archive
rm ~/GoShareShares/aB3x9K2m.json ~/GoShareShares/aB3x9K2m.zip
```

## Security Considerations

### What's Protected

✅ **Path traversal** - ZIP extraction validates all paths
✅ **Expiration** - Shares auto-expire after 7 days
✅ **Local only** - No cloud storage unless you upload the ZIP
✅ **No metadata leakage** - `.goshare.json` excluded from archive

### What's NOT Protected

❌ **Encryption** - ZIP is not encrypted
❌ **Access control** - Anyone with code can download
❌ **Audit trail** - No log of who downloaded
❌ **Revocation** - Can't invalidate a code early (delete files manually)

### Best Practices

1. **Use secure channels** for sharing codes
2. **Set expectations** - "Code expires in 7 days"
3. **Don't share secrets** - No API keys, passwords, tokens
4. **Clean up** - Delete shares after confirmed download
5. **Verify** - Ask recipient to confirm receipt

## Troubleshooting Shares

### "Share code not found or expired"

```bash
# Check if files exist
ls -la ~/GoShareShares/aB3x9K2m.*

# Check expiration
cat ~/GoShareShares/aB3x9K2m.json | jq '.expires_at'
# Compare with current time: date -u +"%Y-%m-%dT%H:%M:%SZ"
```

### "Archive not found for share code"

```bash
# ZIP might have been deleted manually
# Regenerate share:
goshare share -project my-project
```

### "Invalid share info"

```bash
# JSON corrupted
rm ~/GoShareShares/aB3x9K2m.json ~/GoShareShares/aB3x9K2m.zip
goshare share -project my-project
```

### Download extracts to wrong location

```bash
# Use -setup flag
goshare download -share aB3x9K2m -setup /exact/path/you/want
```

## Share Code Format Details

```
Characters: A-Z, a-z, 0-9, -, _ (URL-safe base64)
Length: 8 characters
Entropy: 6 bytes = 48 bits
Collision probability: ~1 in 2.8 trillion
Generation: crypto/rand (cryptographically secure)
```

## Integration Examples

### Share via GitHub Actions

```yaml
# .github/workflows/share.yml
name: Share Release
on:
  release:
    types: [published]
jobs:
  share:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: make build
      - run: |
          ./bin/goshare create -project ${{ github.event.repository.name }} \
            -summary "Release ${{ github.event.release.tag_name }}"
      - run: ./bin/goshare cp -project ${{ github.event.repository.name }} -cp ./dist
      - run: |
          CODE=$(./bin/goshare share -project ${{ github.event.repository.name }} 2>&1 | grep "Share code:" | awk '{print $3}')
          echo "SHARE_CODE=$CODE" >> $GITHUB_ENV
      - name: Comment on release
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: ${{ github.event.release.id }},
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `GoShare download code: \`${process.env.SHARE_CODE}\`\n\nDownload with: \`goshare download -share ${process.env.SHARE_CODE}\``
            })
```

### Share via Script

```bash
#!/bin/bash
# quick-share.sh

PROJECT=$1
if [ -z "$PROJECT" ]; then
  echo "Usage: $0 <project-name>"
  exit 1
fi

CODE=$(goshare share -project "$PROJECT" 2>&1 | grep "Share code:" | awk '{print $3}')

if [ -z "$CODE" ]; then
  echo "Failed to create share"
  exit 1
fi

echo "✅ Share created for '$PROJECT'"
echo "📋 Code: $CODE"
echo "📥 Download: goshare download -share $CODE"
echo "⏰ Expires: $(date -d '+7 days' '+%Y-%m-%d %H:%M')"

# Copy to clipboard (Linux)
echo -n "$CODE" | xclip -selection clipboard 2>/dev/null || true
# macOS: echo -n "$CODE" | pbcopy
# Windows: echo -n "$CODE" | clip
```

## FAQ

**Q: Can I extend a share's expiration?**
A: No, but you can regenerate: `goshare share -project my-project`

**Q: Can I password-protect a share?**
A: Not built-in. Encrypt the ZIP manually before sharing.

**Q: How many times can a code be used?**
A: Unlimited, until expiration.

**Q: Can I share a project without the index.html?**
A: The index.html is always included. Delete it from the project first if needed.

**Q: Does the recipient need GoShare installed?**
A: Yes, to use `goshare download`. They could manually extract the ZIP if you send it directly.