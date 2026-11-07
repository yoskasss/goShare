package publish

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

type Publisher struct{}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(projectName string) error {
	homeDir, _ := os.UserHomeDir()
	projectPath := filepath.Join(homeDir, "GoShareProjects", projectName)

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project '%s' not found", projectName)
	}

	indexPath := filepath.Join(projectPath, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return fmt.Errorf("index.html not found in project")
	}

	color.Cyan("Starting cloudflared tunnel for project '%s'...", projectName)
	color.Yellow("Press Ctrl+C to stop the tunnel")

	cmd := exec.Command("cloudflared", "tunnel", "--url", fmt.Sprintf("http://localhost:8080"))
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		time.Sleep(2 * time.Second)
		color.Green("Tunnel started! Share the URL shown above.")
		color.Cyan("Serving files from: %s", projectPath)
	}()

	if err := p.serveHTTP(projectPath); err != nil {
		return err
	}

	return cmd.Wait()
}

func (p *Publisher) serveHTTP(projectPath string) error {
	http.Handle("/", http.FileServer(http.Dir(projectPath)))
	color.Cyan("Starting local server on :8080...")
	return http.ListenAndServe(":8080", nil)
}

func (p *Publisher) CheckCloudflared() error {
	_, err := exec.LookPath("cloudflared")
	if err != nil {
		return fmt.Errorf("cloudflared not found in PATH. Install it from https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation")
	}
	return nil
}
