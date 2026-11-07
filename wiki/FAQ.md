# FAQ

Frequently Asked Questions about GoShare.

## General

### What is GoShare?

GoShare is a lightweight command-line tool for creating, managing, and sharing Go (or any) project directories. It generates a GitHub-style HTML project page, lets you copy files into projects, share them via short codes, and publish them publicly via Cloudflare tunnels.

### Who is it for?

- Developers who want to quickly share project snapshots
- Teams needing lightweight project sharing without Git
- Anyone wanting to publish a project directory to the web temporarily
- Teachers/students sharing code assignments

### Is it only for Go projects?

**No!** Despite the name, GoShare works with any project type:
- Go, Python, JavaScript, Rust, etc.
- Documentation projects
- Static websites
- Any directory structure

The name "GoShare" comes from being written in Go, not for Go projects only.

### How is it different from Git/GitHub?

| Feature | GoShare | Git/GitHub |
|---------|---------|------------|
| Version control | ❌ No | ✅ Yes |
| Collaboration | ❌ No | ✅ Yes |
| Quick sharing | ✅ 8-char code | ❌ Requires repo setup |
| Web publishing | ✅ One command | ⚠️ Pages/Action setup |
| History | ❌ No | ✅ Full history |
| Authentication | ❌ No | ✅ Built-in |
| File size limits | Disk only | 100MB/file, 2GB repo |

**Use GoShare for:** Quick sharing, demos, temporary publishing, code reviews without Git
**Use Git for:** Version control, collaboration, long-term projects

### Is there a web interface?

No, GoShare is CLI-only. The generated `index.html` provides a web-viewable project page, but all management is via command line.

### Does it require a server?

**No.** GoShare runs entirely locally. Data is stored in your home directory. Only the `public` command creates a temporary public tunnel via Cloudflare.

## Installation

### Do I need Go installed to use it?

**No**, if you download a pre-built binary. **Yes**, if building from source.

### Can I install it via package manager?

Not yet. Currently: download binary or build from source. Homebrew/Scoop/AUR packages planned.

### Where are the binaries?

GitHub Releases page (when published). For now, build from source: `make build-all`

### Does it work on Windows?

Yes! Build with `GOOS=windows GOARCH=amd64 go build` or use `make build-all`.

## Projects

### Where are projects stored?

`~/GoShareProjects/<project-name>/`

### Can I change the storage location?

Not via config yet. Modify `internal/project/manager.go` and rebuild.

### What's in a project directory?

```
my-project/
├── .goshare.json    # Metadata (name, summary, timestamps)
├── index.html       # Generated project page
└── your-files/      # Whatever you copied in
```

### Can I have nested projects?

No, projects are flat in `~/GoShareProjects/`. But you can have any directory structure *inside* a project.

### How do I rename a project?

```bash
# No direct rename command
# Workaround:
goshare create -project new-name -summary "Same project"
goshare cp -project new-name -cp ~/GoShareProjects/old-name
goshare delete -project old-name -force
```

### Can I move a project to another machine?

Yes! Copy `~/GoShareProjects/my-project/` to the new machine's `~/GoShareProjects/`. Run `goshare list` to verify.

### What happens if I delete the .goshare.json file?

The project becomes "orphaned" - files remain but GoShare won't recognize it. You can recreate metadata by running `goshare create` with the same name.

## Sharing

### How long do share codes last?

**7 days** from creation. Then they're automatically cleaned up.

### Can I extend expiration?

Not directly. Regenerate: `goshare share -project my-project`

### Can I make a share permanent?

No. For permanent sharing, use Git/GitHub or copy the ZIP manually.

### How many times can a code be used?

Unlimited, until expiration.

### Can I revoke a share early?

Yes, delete the files:
```bash
rm ~/GoShareShares/<code>.json ~/GoShareShares/<code>.zip
```

### Is the share encrypted?

**No.** The ZIP is not encrypted. Anyone with the code can download. Share codes via secure channels.

### Can I password-protect a share?

Not built-in. Workaround: encrypt the ZIP yourself before sharing the code.

### What's the maximum project size for sharing?

Limited by disk space. The ZIP is created in `~/GoShareShares/`. Large projects = large ZIPs = slow sharing.

### Can I share just one file?

Share the whole project, or copy just that file to a new project first:
```bash
goshare create -project single-file -summary "Just one file"
goshare cp -project single-file -cp ~/path/to/file.txt
goshare share -project single-file
```

### Does the recipient need GoShare?

Yes, for `goshare download`. But you could send them the ZIP directly (from `~/GoShareShares/<code>.zip`) and they can extract manually.

## Publishing

### What is "publishing"?

The `public` command starts a local web server + Cloudflare tunnel, giving you a public HTTPS URL like `https://abc-123.trycloudflare.com` that serves your project directory.

### Do I need a Cloudflare account?

**No** for quick tunnels (what `public` uses). **Yes** for named tunnels with custom domains.

### How long does the tunnel last?

Until you press `Ctrl+C`. Quick tunnels are temporary.

### Can I use a custom domain?

Not with `public` command. Set up a named Cloudflare Tunnel separately (see Publishing wiki).

### Is it secure?

⚠️ **Everything in the project is publicly accessible.** Review for secrets before publishing. Stop tunnel when done.

### Can multiple people view simultaneously?

Yes, Cloudflare handles concurrent connections.

### Does it work behind corporate firewall?

Usually yes - Cloudflare uses standard HTTPS (port 443).

### What's the bandwidth limit?

Free tier is generous. Check Cloudflare pricing for details.

## Technical

### What Go version is required?

Go 1.25+ for building. Binary runs on any compatible OS/arch.

### What dependencies does it have?

- `github.com/fatih/color` - Terminal colors
- `github.com/common-nighthawk/go-figure` - ASCII banner
- `github.com/n-marshall/go-cp` - File copying
- Standard library for everything else

### Is it open source?

Yes, MIT License.

### Can I contribute?

Yes! See Development wiki. PRs welcome.

### How do I report a bug?

GitHub Issues with:
- OS/version
- GoShare version (`goshare version`)
- Command that failed
- Full error output
- Steps to reproduce

### Is there an API?

No, CLI only. But you can call the Go packages directly in your own code.

### Can I use it in CI/CD?

Yes! Build the binary, then run commands in your pipeline. See Development wiki for examples.

## Troubleshooting

### "command not found"

Binary not in PATH. See Installation wiki.

### "project already exists"

Use different name or delete old: `goshare delete -project name -force`

### "share code not found"

Code expired (7 days) or typo. Regenerate with `goshare share`.

### "cloudflared not found"

Install cloudflared. See Publishing wiki.

### "port 8080 in use"

Kill process: `lsof -i :8080` then `kill -9 <PID>`

### Permission denied

Check `~/GoShareProjects/` permissions: `ls -ld ~/GoShareProjects/`

## Roadmap / Planned Features

- [ ] Config file (`~/.config/goshare.yaml`)
- [ ] Password-protected shares
- [ ] Custom share expiration
- [ ] Project rename command
- [ ] Share listing command
- [ ] Web UI (optional)
- [ ] Git integration (optional)
- [ ] Plugin system
- [ ] Package manager distribution

## Still Have Questions?

- Check other wiki pages
- Search GitHub Issues
- Open a new Issue
- Check the source code (it's readable!)