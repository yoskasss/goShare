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
	case "git":
		return a.runGit(args[2:])
	case "collab":
		return a.runCollab(args[2:])
	case "history":
		return a.runHistory(args[2:])
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
	fmt.Println("  git        Git version control commands")
	fmt.Println("  collab     Collaboration management")
	fmt.Println("  history    Project history and commits")
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
	fmt.Println()
	color.Yellow("Git Commands:")
	fmt.Println("  goshare git init -project myproject")
	fmt.Println("  goshare git status -project myproject")
	fmt.Println("  goshare git add -project myproject -files file1.go,file2.go")
	fmt.Println("  goshare git commit -project myproject -msg \"commit message\"")
	fmt.Println("  goshare git log -project myproject -limit 10")
	fmt.Println("  goshare git diff -project myproject")
	fmt.Println("  goshare git branch -project myproject")
	fmt.Println("  goshare git checkout -project myproject -branch main")
	fmt.Println("  goshare git new-branch -project myproject -branch feature")
	fmt.Println()
	color.Yellow("Collaboration Commands:")
	fmt.Println("  goshare collab add -project myproject -email user@example.com")
	fmt.Println("  goshare collab remove -project myproject -email user@example.com")
	fmt.Println("  goshare collab list -project myproject")
	fmt.Println()
	color.Yellow("History Commands:")
	fmt.Println("  goshare history log -project myproject -limit 10")
	fmt.Println("  goshare history show -project myproject -hash abc123")
	fmt.Println("  goshare history export -project myproject -format markdown")
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

// ============ GIT COMMANDS ============

func (a *App) runGit(args []string) error {
	if len(args) < 1 {
		color.Red("Git subcommand required")
		fmt.Println("Usage: goshare git <init|status|add|commit|log|diff|branch|checkout|new-branch> [options]")
		return nil
	}

	subcommand := args[0]

	switch subcommand {
	case "init":
		return a.runGitInit(args[1:])
	case "status":
		return a.runGitStatus(args[1:])
	case "add":
		return a.runGitAdd(args[1:])
	case "commit":
		return a.runGitCommit(args[1:])
	case "log":
		return a.runGitLog(args[1:])
	case "diff":
		return a.runGitDiff(args[1:])
	case "branch":
		return a.runGitBranch(args[1:])
	case "checkout":
		return a.runGitCheckout(args[1:])
	case "new-branch":
		return a.runGitNewBranch(args[1:])
	default:
		color.Red("Unknown git subcommand: %s", subcommand)
		fmt.Println("Available: init, status, add, commit, log, diff, branch, checkout, new-branch")
		return fmt.Errorf("unknown git subcommand: %s", subcommand)
	}
}

func (a *App) runGitInit(args []string) error {
	fs := flag.NewFlagSet("git init", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	color.Cyan("Initializing git repository for project: %s", *projectName)
	return a.projectManager.GitInit(*projectName)
}

func (a *App) runGitStatus(args []string) error {
	fs := flag.NewFlagSet("git status", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	status, err := a.projectManager.GitStatus(*projectName)
	if err != nil {
		return err
	}

	color.Cyan("Git Status for %s:", *projectName)
	fmt.Println(status)
	return nil
}

func (a *App) runGitAdd(args []string) error {
	fs := flag.NewFlagSet("git add", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	filesStr := fs.String("files", "", "Comma-separated list of files to add (or '.' for all)")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}
	if *filesStr == "" {
		return fmt.Errorf("files list is required (-files)")
	}

	files := strings.Split(*filesStr, ",")
	for i := range files {
		files[i] = strings.TrimSpace(files[i])
	}

	color.Cyan("Adding files to git: %v", files)
	return a.projectManager.GitAdd(*projectName, files)
}

func (a *App) runGitCommit(args []string) error {
	fs := flag.NewFlagSet("git commit", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	message := fs.String("msg", "", "Commit message")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}
	if *message == "" {
		return fmt.Errorf("commit message is required (-msg)")
	}

	color.Cyan("Committing changes: %s", *message)
	return a.projectManager.GitCommit(*projectName, *message)
}

func (a *App) runGitLog(args []string) error {
	fs := flag.NewFlagSet("git log", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	limit := fs.Int("limit", 10, "Number of commits to show")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	commits, err := a.projectManager.GitLog(*projectName, *limit)
	if err != nil {
		return err
	}

	color.Cyan("Git Log for %s (last %d commits):", *projectName, *limit)
	fmt.Println()
	for _, c := range commits {
		color.Yellow("  %s", c.Hash)
		fmt.Printf("    Author: %s <%s>\n", c.Author, c.Email)
		fmt.Printf("    Date:   %s\n", c.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Message: %s\n", c.Message)
		if len(c.Files) > 0 {
			fmt.Printf("    Files: %s\n", strings.Join(c.Files, ", "))
		}
		fmt.Println()
	}
	return nil
}

func (a *App) runGitDiff(args []string) error {
	fs := flag.NewFlagSet("git diff", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	file := fs.String("file", "", "Specific file to diff")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	diff, err := a.projectManager.GitDiff(*projectName, *file)
	if err != nil {
		return err
	}

	color.Cyan("Git Diff for %s:", *projectName)
	fmt.Println(diff)
	return nil
}

func (a *App) runGitBranch(args []string) error {
	fs := flag.NewFlagSet("git branch", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	branches, err := a.projectManager.GitBranch(*projectName)
	if err != nil {
		return err
	}

	color.Cyan("Branches for %s:", *projectName)
	for _, b := range branches {
		fmt.Printf("  %s\n", b)
	}
	return nil
}

func (a *App) runGitCheckout(args []string) error {
	fs := flag.NewFlagSet("git checkout", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	branch := fs.String("branch", "", "Branch name to checkout")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}
	if *branch == "" {
		return fmt.Errorf("branch name is required (-branch)")
	}

	color.Cyan("Checking out branch: %s", *branch)
	return a.projectManager.GitCheckout(*projectName, *branch)
}

func (a *App) runGitNewBranch(args []string) error {
	fs := flag.NewFlagSet("git new-branch", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	branch := fs.String("branch", "", "New branch name")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}
	if *branch == "" {
		return fmt.Errorf("branch name is required (-branch)")
	}

	color.Cyan("Creating and switching to new branch: %s", *branch)
	return a.projectManager.GitCreateBranch(*projectName, *branch)
}

// ============ COLLABORATION COMMANDS ============

func (a *App) runCollab(args []string) error {
	if len(args) < 1 {
		color.Red("Collab subcommand required")
		fmt.Println("Usage: goshare collab <add|remove|list> [options]")
		return nil
	}

	subcommand := args[0]

	switch subcommand {
	case "add":
		return a.runCollabAdd(args[1:])
	case "remove":
		return a.runCollabRemove(args[1:])
	case "list":
		return a.runCollabList(args[1:])
	default:
		color.Red("Unknown collab subcommand: %s", subcommand)
		fmt.Println("Available: add, remove, list")
		return fmt.Errorf("unknown collab subcommand: %s", subcommand)
	}
}

func (a *App) runCollabAdd(args []string) error {
	fs := flag.NewFlagSet("collab add", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	email := fs.String("email", "", "Collaborator email")
	name := fs.String("name", "", "Collaborator name")
	role := fs.String("role", "contributor", "Role: owner, contributor, viewer")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}
	if *email == "" {
		return fmt.Errorf("collaborator email is required (-email)")
	}

	color.Cyan("Adding collaborator %s (%s) to project %s", *email, *role, *projectName)
	return a.projectManager.AddCollaborator(*projectName, *name, *email, *role)
}

func (a *App) runCollabRemove(args []string) error {
	fs := flag.NewFlagSet("collab remove", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	email := fs.String("email", "", "Collaborator email to remove")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}
	if *email == "" {
		return fmt.Errorf("collaborator email is required (-email)")
	}

	color.Cyan("Removing collaborator %s from project %s", *email, *projectName)
	return a.projectManager.RemoveCollaborator(*projectName, *email)
}

func (a *App) runCollabList(args []string) error {
	fs := flag.NewFlagSet("collab list", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	collaborators, err := a.projectManager.ListCollaborators(*projectName)
	if err != nil {
		return err
	}

	color.Cyan("Collaborators for %s:", *projectName)
	if len(collaborators) == 0 {
		fmt.Println("  No collaborators")
		return nil
	}
	for _, c := range collaborators {
		fmt.Printf("  %s\n", c)
	}
	return nil
}

// ============ HISTORY COMMANDS ============

func (a *App) runHistory(args []string) error {
	if len(args) < 1 {
		color.Red("History subcommand required")
		fmt.Println("Usage: goshare history <log|show|export> [options]")
		return nil
	}

	subcommand := args[0]

	switch subcommand {
	case "log":
		return a.runHistoryLog(args[1:])
	case "show":
		return a.runHistoryShow(args[1:])
	case "export":
		return a.runHistoryExport(args[1:])
	default:
		color.Red("Unknown history subcommand: %s", subcommand)
		fmt.Println("Available: log, show, export")
		return fmt.Errorf("unknown history subcommand: %s", subcommand)
	}
}

func (a *App) runHistoryLog(args []string) error {
	fs := flag.NewFlagSet("history log", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	limit := fs.Int("limit", 10, "Number of commits to show")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	commits, err := a.projectManager.GetHistory(*projectName, *limit)
	if err != nil {
		return err
	}

	color.Cyan("History for %s (last %d commits):", *projectName, *limit)
	fmt.Println()
	for _, c := range commits {
		color.Yellow("  %s", c.Hash)
		fmt.Printf("    Author: %s <%s>\n", c.Author, c.Email)
		fmt.Printf("    Date:   %s\n", c.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Message: %s\n", c.Message)
		if len(c.Files) > 0 {
			fmt.Printf("    Files: %s\n", strings.Join(c.Files, ", "))
		}
		fmt.Println()
	}
	return nil
}

func (a *App) runHistoryShow(args []string) error {
	fs := flag.NewFlagSet("history show", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	hash := fs.String("hash", "", "Commit hash")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}
	if *hash == "" {
		return fmt.Errorf("commit hash is required (-hash)")
	}

	commit, err := a.projectManager.ShowCommit(*projectName, *hash)
	if err != nil {
		return err
	}

	color.Cyan("Commit Details:")
	color.Yellow("  Hash: %s", commit.Hash)
	fmt.Printf("  Author: %s <%s>\n", commit.Author, commit.Email)
	fmt.Printf("  Date:   %s\n", commit.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Message: %s\n", commit.Message)
	if len(commit.Files) > 0 {
		fmt.Printf("  Files:\n")
		for _, f := range commit.Files {
			fmt.Printf("    - %s\n", f)
		}
	}
	return nil
}

func (a *App) runHistoryExport(args []string) error {
	fs := flag.NewFlagSet("history export", flag.ExitOnError)
	projectName := fs.String("project", "", "Project name")
	format := fs.String("format", "text", "Export format: text, json, markdown")
	fs.Parse(args)

	if *projectName == "" {
		return fmt.Errorf("project name is required (-project)")
	}

	output, err := a.projectManager.ExportHistory(*projectName, *format)
	if err != nil {
		return err
	}

	fmt.Println(output)
	return nil
}
