# Quick Start

Get up and running with GoShare in 5 minutes.

## 1. Create Your First Project

```bash
# Interactive (prompts for summary)
goshare create

# Or with flags
goshare create -project my-awesome-project -summary "A fantastic Go project"
```

Output:
```
Project 'my-awesome-project' created successfully at /home/user/GoShareProjects/my-awesome-project
```

## 2. Add Files to Your Project

```bash
# Copy a directory
goshare copy -project my-awesome-project -cp ./src

# Copy a single file
goshare copy -project my-awesome-project -cp ./README.md

# Shorthand
goshare cp -project my-awesome-project -cp ./src
```

Output:
```
Copying ./src to /home/user/GoShareProjects/my-awesome-project/src...
Copied successfully!
```

## 3. View Your Project

```bash
# Open the generated HTML page
xdg-open ~/GoShareProjects/my-awesome-project/index.html
# Or on macOS:
open ~/GoShareProjects/my-awesome-project/index.html
```

The page shows:
- Project name and summary
- Creation/update timestamps
- File tree with all project files

## 4. Share Your Project

```bash
# Generate a share code
goshare share -project my-awesome-project
```

Output:
```
Creating archive for sharing...
Share code: aB3x9K2m
Download with: goshare download -share aB3x9K2m
```

**Share the code `aB3x9K2m` with anyone.** They can download it without needing your original files.

## 5. Download a Shared Project

```bash
# Receiver runs this
goshare download -share aB3x9K2m

# With custom destination
goshare download -share aB3x9K2m -setup /path/to/destination
```

Output:
```
Downloading project 'my-awesome-project'...
Extracted successfully to /home/user/GoShareProjects/my-awesome-project
Project downloaded to: /home/user/GoShareProjects/my-awesome-project
```

## 6. Publish to the Web (Optional)

Requires `cloudflared` installed.

```bash
# Start public tunnel
goshare public -project my-awesome-project
```

Output:
```
Starting cloudflared tunnel for project 'my-awesome-project'...
Press Ctrl+C to stop the tunnel
Tunnel started! Share the URL shown above.
Serving files from: /home/user/GoShareProjects/my-awesome-project
```

Cloudflare will display a URL like `https://random-name.trycloudflare.com` - share this with anyone to view your project live.

Press `Ctrl+C` to stop the tunnel.

## 7. List All Projects

```bash
goshare list
# Or shorthand
goshare ls
```

Output:
```
Projects:
  my-awesome-project - A fantastic Go project
  another-project    - Another cool project
```

## 8. Clean Up

```bash
# Delete with confirmation
goshare delete -project my-awesome-project

# Force delete (no confirmation)
goshare delete -project my-awesome-project -force
```

Output:
```
Project 'my-awesome-project' deleted successfully
```

---

## Complete Workflow Example

```bash
# 1. Create project
goshare create -project blog -summary "My personal blog engine"

# 2. Add source code
goshare cp -project blog -cp ~/projects/blog-engine

# 3. Share with colleague
goshare share -project blog
# Share code: Xy7zP2qR

# 4. Colleague downloads
goshare download -share Xy7zP2qR

# 5. Publish for demo
goshare public -project blog
# Share the cloudflare URL

# 6. List projects
goshare ls

# 7. Clean up when done
goshare delete -project blog -force
```

## Pro Tips

| Tip | Command |
|-----|---------|
| Tab completion | `goshare <TAB>` (after setup) |
| Short flags | `-p` for `-project`, `-c` for `-cp` |
| View help | `goshare help` or `goshare create -h` |
| Version check | `goshare version` |

## Next Steps

- [Commands Reference](Commands-Reference.md) - Detailed command documentation
- [Sharing Projects](Sharing-Projects.md) - Advanced sharing workflows
- [Publishing with Cloudflare](Publishing-with-Cloudflare.md) - Tunnel configuration
- [Project Structure](Project-Structure.md) - Understanding the file layout