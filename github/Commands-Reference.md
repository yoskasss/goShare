# Commands Reference

Complete reference for all GoShare commands.

## Global Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `-project` | `-p` | Project name | `GoShare` |
| `-cp` | `-c` | Source file/directory to copy | - |
| `-share` | `-s` | Share code for download | - |
| `-setup` | - | Custom destination for download | - |
| `-force` | `-f` | Force delete without confirmation | `false` |
| `-help` | `-h` | Show help for command | - |

## Commands

### `create` - Create a New Project

```bash
goshare create [flags]
```

**Flags:**
- `-project`, `-p` - Project name (default: "GoShare")
- `-summary` - Project summary (optional, prompts if omitted)

**Examples:**
```bash
# Interactive
goshare create

# With flags
goshare create -project myapp -summary "My Go application"

# Short flags
goshare create -p myapp -summary "My Go application"
```

**What happens:**
1. Creates `~/GoShareProjects/<project>/`
2. Creates `.goshare.json` with metadata
3. Generates `index.html` (GitHub-style project page)
4. Opens in browser if `$BROWSER` is set

---

### `copy` / `cp` - Copy Files into Project

```bash
goshare copy [flags]
goshare cp [flags]          # Shorthand
```

**Flags:**
- `-project`, `-p` - Target project name **(required)**
- `-cp`, `-c` - Source file or directory **(required)**

**Examples:**
```bash
# Copy directory
goshare copy -project myapp -cp ./src

# Copy single file
goshare cp -p myapp -c ./README.md

# Copy from absolute path
goshare copy -p myapp -cp /home/user/projects/mylib
```

**Behavior:**
- Preserves directory structure
- Updates project's `index.html` file tree
- Updates `UpdatedAt` timestamp in metadata
- Skips hidden files (except `.goshare.json`)

---

### `public` - Publish Project via Cloudflare Tunnel

```bash
goshare public [flags]
```

**Flags:**
- `-project`, `-p` - Project name to publish **(required)**

**Prerequisites:**
- `cloudflared` installed and in PATH
- Project must exist with `index.html`

**Examples:**
```bash
goshare public -project myapp
goshare public -p myapp
```

**What happens:**
1. Starts local HTTP server on port 8080
2. Creates Cloudflare quick tunnel
3. Displays public URL (e.g., `https://abc-123.trycloudflare.com`)
4. Serves entire project directory
5. Press `Ctrl+C` to stop

**Output:**
```
Starting cloudflared tunnel for project 'myapp'...
Press Ctrl+C to stop the tunnel
Tunnel started! Share the URL shown above.
Serving files from: /home/user/GoShareProjects/myapp
```

---

### `share` - Generate Share Code

```bash
goshare share [flags]
```

**Flags:**
- `-project`, `-p` - Project name to share **(required)**

**Examples:**
```bash
goshare share -project myapp
goshare share -p myapp
```

**Output:**
```
Creating archive for sharing...
Share code: aB3x9K2m
Download with: goshare download -share aB3x9K2m
```

**What happens:**
1. Creates ZIP archive of project (excludes `.goshare.json`)
2. Stores in `~/GoShareShares/<code>.zip`
3. Creates metadata file `~/GoShareShares/<code>.json`
4. Generates 8-character share code
5. Share expires in **7 days**

---

### `download` - Download Shared Project

```bash
goshare download [flags]
```

**Flags:**
- `-share`, `-s` - Share code **(required)**
- `-setup` - Custom destination directory (optional)

**Examples:**
```bash
# Default location (~/GoShareProjects/<project-name>)
goshare download -share aB3x9K2m

# Custom destination
goshare download -share aB3x9K2m -setup /tmp/my-download
goshare download -s aB3x9K2m -setup ~/projects/restored
```

**Output:**
```
Downloading project 'myapp'...
Extracted successfully to /home/user/GoShareProjects/myapp
Project downloaded to: /home/user/GoShareProjects/myapp
```

**Behavior:**
- Validates share code exists and not expired
- Extracts ZIP to destination
- Creates project directory with original name
- Overwrites if directory exists

---

### `delete` - Delete a Project

```bash
goshare delete [flags]
```

**Flags:**
- `-project`, `-p` - Project name to delete **(required)**
- `-force`, `-f` - Skip confirmation prompt

**Examples:**
```bash
# With confirmation
goshare delete -project myapp

# Force delete
goshare delete -project myapp -force
goshare delete -p myapp -f
```

**Output (with confirmation):**
```
Are you sure you want to delete project 'myapp'? (y/N): y
Project 'myapp' deleted successfully
```

**Output (force):**
```
Project 'myapp' deleted successfully
```

**Warning:** This permanently deletes the project directory and all files.

---

### `list` / `ls` - List All Projects

```bash
goshare list
goshare ls              # Shorthand
```

**Examples:**
```bash
goshare list
goshare ls
```

**Output:**
```
Projects:
  myapp           - My Go application
  blog-engine     - Personal blog engine
  cli-tool        - Command line utility
```

**Shows:** Project name and summary from metadata.

---

### `help` - Show Help

```bash
goshare help
goshare -h
goshare --help
```

**Output:** Shows all commands with descriptions and examples.

---

### `version` - Show Version

```bash
goshare version
goshare -v
goshare --version
```

**Output:**
```
GoShare v1.0.0
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (invalid args, not found, etc.) |
| 2 | Usage error (wrong flags) |

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `HOME` | User home directory (used for data paths) |
| `BROWSER` | Browser command to open index.html (optional) |

---

## Configuration Files

### Project Metadata (`~/.goshare.json`)

```json
{
  "name": "myapp",
  "summary": "My Go application",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:45:00Z",
  "path": "/home/user/GoShareProjects/myapp",
  "index_html": "/home/user/GoShareProjects/myapp/index.html"
}
```

### Share Metadata (`~/GoShareShares/<code>.json`)

```json
{
  "code": "aB3x9K2m",
  "project_name": "myapp",
  "created_at": "2024-01-15T12:00:00Z",
  "expires_at": "2024-01-22T12:00:00Z",
  "zip_path": "/home/user/GoShareShares/aB3x9K2m.zip"
}
```

---

## Shell Completion

```bash
# Generate for your shell
goshare completion bash
goshare completion zsh
goshare completion fish
goshare completion powershell
```

See [Installation](Installation.md#shell-completion) for setup instructions.