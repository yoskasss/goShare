package cli

import (
	"flag"
	"fmt"
	"strings"

	"GoShare/internal/project"
	"GoShare/internal/publish"
	"GoShare/internal/share"
	"GoShare/internal/ziputil"

	"github.com/fatih/color"
)

type App struct {
	projectManager *project.Manager
	publisher      *publish.Publisher
	shareManager   *share.Manager
	zipManager     *ziputil.Manager
}

func NewApp() *App {
	return &App{
		projectManager: project.NewManager(),
		publisher:      publish.NewPublisher(),
		shareManager:   share.NewManager(),
		zipManager:     ziputil.NewManager(),
	}
}

func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.printUsage()
		return nil
	}

	command := args[1]

	switch command {
	case "create":
		return a.runCreate(args[2:])
	case "copy", "cp":
		return a.runCopy(args[2:])
	case "public":
		return a.runPublic(args[2:])
	case "share":
		return a.runShare(args[2:])
	case "download":
		return a.runDownload(args[2:])
	case "delete":
		return a.runDelete(args[2:])
	case "list", "ls":
		return a.runList(args[2:])
	case "help", "-h", "--help":
		a.printUsage()
		return nil
	case "version", "-v", "--version":
		fmt.Println("GoShare v1.0.0")
		return nil
	default:
		color.Red("Unknown command: %s", command)
		a.printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (a *App) printUsage() {
	color.Cyan("GoShare - Project sharing tool")
	fmt.Println()
	color.Yellow("Usage:")
	fmt.Println("  goshare <command> [options]")
	fmt.Println()
	color.Yellow("Commands:")
	fmt.Println("  create     Create a new project")
	fmt.Println("  copy, cp   Copy files/directories into a project")
	fmt.Println("  public     Publish a project with cloudflared")
	fmt.Println("  share      Generate a share code for a project")
	fmt.Println("  download   Download a shared project")
	fmt.Println("  delete     Delete a project")
	fmt.Println("  list, ls   List all projects")
	fmt.Println("  help       Show this help")
	fmt.Println("  version    Show version")
	fmt.Println()
	color.Yellow("Options:")
	fmt.Println("  -project   Project name (default: GoShare)")
	fmt.Println("  -cp        Source file/directory to copy")
	fmt.Println("  -public    Project name to publish")
	fmt.Println("  -share     Share code to download")
	fmt.Println("  -setup     Setup archive path")
	fmt.Println("  -delete    Project name to delete")
	fmt.Println()
	color.Yellow("Examples:")
	fmt.Println("  goshare create -project myproject")
	fmt.Println("  goshare copy -project myproject -cp ./src")
	fmt.Println("  goshare public -project myproject")
	fmt.Println("  goshare share -project myproject")
	fmt.Println("  goshare download -share ABC123")
	fmt.Println("  goshare delete -project myproject")
	fmt.Println("  goshare list")
}

func (a *App) runCreate(args []string) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	projectName := fs.String("project", "GoShare", "Project name")
	summary := fs.String("summary", "", "Project summary")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required")
	}

	return a.projectManager.Create(*projectName, *summary)
}

func (a *App) runCopy(args []string) error {
	fs := flag.NewFlagSet("copy", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	copySource := fs.String("cp", "", "Source file/directory to copy")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}
	if *copySource == "" {
		return fmt.Errorf("source path is required (-cp)")
	}

	return a.projectManager.Copy(*projectName, *copySource)
}

func (a *App) runPublic(args []string) error {
	fs := flag.NewFlagSet("public", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name to publish")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	return a.publisher.Publish(*projectName)
}

func (a *App) runShare(args []string) error {
	fs := flag.NewFlagSet("share", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name to share")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	code, err := a.shareManager.GenerateShareCode(*projectName)
	if err != nil {
		return err
	}

	color.Green("Share code: %s", code)
	fmt.Printf("Download with: goshare download -share %s\n", code)
	return nil
}

func (a *App) runDownload(args []string) error {
	fs := flag.NewFlagSet("download", flag.ExitOnError)
	shareCode := fs.String("share", "", "Share code")
	setupPath := fs.String("setup", "", "Setup archive path")
	fs.Parse(args)

	if *shareCode == "" {
		return fmt.Errorf("share code is required (-share)")
	}

	return a.shareManager.Download(*shareCode, *setupPath)
}

func (a *App) runDelete(args []string) error {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name to delete")
	force := fs.Bool("force", false, "Force delete without confirmation")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	if !*force {
		fmt.Printf("Are you sure you want to delete project '%s'? (y/N): ", *projectName)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	return a.projectManager.Delete(*projectName)
}

func (a *App) runList(args []string) error {
	projects, err := a.projectManager.List()
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		color.Yellow("No projects found")
		return nil
	}

	color.Cyan("Projects:")
	for _, p := range projects {
		fmt.Printf("  %s - %s\n", color.GreenString(p.Name), p.Summary)
	}
	return nil
}
