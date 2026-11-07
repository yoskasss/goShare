# GoShare

GoShare is a lightweight Go command-line tool for creating, managing, previewing, and sharing project directories. It can generate a repository-style HTML page for a project, copy files or directories into it, publish it through a Cloudflare quick tunnel, and let another user download the project as a ZIP archive.

## Features

- Create project directories interactively
- Collect a project name and summary
- Generate a GitHub-style `index.html` project page
- Display project files, source code, and images in the generated page
- Copy files or complete directories into an existing project
- Publish a project with `cloudflared`
- Generate a share code for project downloads
- Download shared projects as ZIP files
- Extract downloaded ZIP files safely
- Delete project directories when explicitly requested

## Requirements

- Go 1.25 or newer
- `cloudflared` (for publishing) - [Install guide](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation)

## Installation

### From Source

```bash
git clone https://github.com/yourusername/GoShare
cd GoShare
make build
# Binary will be at bin/goshare
```

### Install to PATH

```bash
make install
```

## Usage

### Create a new project

```bash
goshare create -project myproject
# Or interactively:
goshare create
```

### Copy files/directories into a project

```bash
goshare copy -project myproject -cp ./src
# Or shorthand:
goshare cp -project myproject -cp ./src
```

### Publish a project (requires cloudflared)

```bash
goshare public -project myproject
```

This starts a local HTTP server on port 8080 and creates a Cloudflare tunnel to make it publicly accessible.

### Generate a share code

```bash
goshare share -project myproject
```

Outputs a share code like `aB3x9K2m`. Share this code with others.

### Download a shared project

```bash
goshare download -share aB3x9K2m
# With custom destination:
goshare download -share aB3x9K2m -setup /path/to/destination
```

### Delete a project

```bash
goshare delete -project myproject
# Force delete without confirmation:
goshare delete -project myproject -force
```

### List all projects

```bash
goshare list
# Or shorthand:
goshare ls
```

### Help and version

```bash
goshare help
goshare version
```

## Project Structure

```
GoShare/
├── cmd/
│   └── goshare/
│       └── main.go           # Entry point
├── internal/
│   ├── cli/
│   │   └── app.go            # CLI application logic
│   ├── project/
│   │   ├── manager.go        # Project management
│   │   └── manager_test.go   # Tests
│   ├── publish/
│   │   └── publisher.go      # Cloudflare publishing
│   ├── share/
│   │   ├── manager.go        # Share code management
│   │   └── manager_test.go   # Tests
│   └── ziputil/
│       ├── manager.go        # ZIP utilities
│       └── manager_test.go   # Tests
├── goahed/                   # Legacy web assets
├── bin/                      # Build output
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Development

### Run tests

```bash
make test
```

### Run with coverage

```bash
make test-coverage
```

### Format code

```bash
make fmt
```

### Run linter

```bash
make lint
```

### Build for all platforms

```bash
make build-all
```

### Run in development mode

```bash
make dev ARGS="create -project test"
```

## Data Storage

Projects are stored in `~/GoShareProjects/` with the following structure:

```
~/GoShareProjects/
└── myproject/
    ├── .goshare.json      # Project metadata
    ├── index.html         # Generated project page
    └── ...                # Your project files
```

Share archives are stored in `~/GoShareShares/` and expire after 7 days.

## License

MIT License