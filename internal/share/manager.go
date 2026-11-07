package share

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"GoShare/internal/ziputil"

	"github.com/fatih/color"
)

type ShareInfo struct {
	Code        string    `json:"code"`
	ProjectName string    `json:"project_name"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	ZipPath     string    `json:"zip_path"`
}

type Manager struct {
	shareDir   string
	zipManager *ziputil.Manager
}

func NewManager() *Manager {
	homeDir, _ := os.UserHomeDir()
	return &Manager{
		shareDir:   filepath.Join(homeDir, "GoShareShares"),
		zipManager: ziputil.NewManager(),
	}
}

func (m *Manager) GenerateShareCode(projectName string) (string, error) {
	projectPath := filepath.Join(os.Getenv("HOME"), "GoShareProjects", projectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return "", fmt.Errorf("project '%s' not found", projectName)
	}

	code := m.generateCode()

	zipPath := filepath.Join(m.shareDir, code+".zip")
	if err := os.MkdirAll(m.shareDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create share directory: %w", err)
	}

	color.Cyan("Creating archive for sharing...")
	if err := m.zipManager.CreateZip(projectPath, zipPath); err != nil {
		return "", fmt.Errorf("failed to create zip: %w", err)
	}

	shareInfo := &ShareInfo{
		Code:        code,
		ProjectName: projectName,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		ZipPath:     zipPath,
	}

	infoPath := filepath.Join(m.shareDir, code+".json")
	data, err := json.MarshalIndent(shareInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal share info: %w", err)
	}

	if err := os.WriteFile(infoPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write share info: %w", err)
	}

	return code, nil
}

func (m *Manager) generateCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:8]
}

func (m *Manager) Download(shareCode, setupPath string) error {
	infoPath := filepath.Join(m.shareDir, shareCode+".json")

	data, err := os.ReadFile(infoPath)
	if err != nil {
		return fmt.Errorf("share code '%s' not found or expired", shareCode)
	}

	var shareInfo ShareInfo
	if err := json.Unmarshal(data, &shareInfo); err != nil {
		return fmt.Errorf("invalid share info: %w", err)
	}

	if time.Now().After(shareInfo.ExpiresAt) {
		return fmt.Errorf("share code has expired")
	}

	if _, err := os.Stat(shareInfo.ZipPath); os.IsNotExist(err) {
		return fmt.Errorf("archive not found for share code")
	}

	color.Cyan("Downloading project '%s'...", shareInfo.ProjectName)

	destDir := filepath.Join(os.Getenv("HOME"), "GoShareProjects", shareInfo.ProjectName)
	if setupPath != "" {
		destDir = setupPath
	}

	if err := m.zipManager.ExtractZip(shareInfo.ZipPath, destDir); err != nil {
		return fmt.Errorf("failed to extract zip: %w", err)
	}

	color.Green("Project downloaded to: %s", destDir)
	return nil
}

func (m *Manager) DownloadFromURL(url, destPath string) error {
	color.Cyan("Downloading from %s...", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	tmpFile := filepath.Join(os.TempDir(), "goshare_download_"+m.generateCode()+".zip")
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write download: %w", err)
	}

	if destPath == "" {
		destPath = filepath.Join(os.Getenv("HOME"), "GoShareProjects", "downloaded_"+m.generateCode())
	}

	if err := m.zipManager.ExtractZip(tmpFile, destPath); err != nil {
		return fmt.Errorf("failed to extract: %w", err)
	}

	os.Remove(tmpFile)
	color.Green("Downloaded and extracted to: %s", destPath)
	return nil
}

func (m *Manager) ListShares() ([]*ShareInfo, error) {
	entries, err := os.ReadDir(m.shareDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*ShareInfo{}, nil
		}
		return nil, err
	}

	var shares []*ShareInfo
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			data, err := os.ReadFile(filepath.Join(m.shareDir, entry.Name()))
			if err != nil {
				continue
			}
			var info ShareInfo
			if json.Unmarshal(data, &info) == nil {
				if time.Now().Before(info.ExpiresAt) {
					shares = append(shares, &info)
				}
			}
		}
	}

	return shares, nil
}

func (m *Manager) CleanupExpired() error {
	entries, err := os.ReadDir(m.shareDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			data, _ := os.ReadFile(filepath.Join(m.shareDir, entry.Name()))
			var info ShareInfo
			if json.Unmarshal(data, &info) == nil {
				if time.Now().After(info.ExpiresAt) {
					os.Remove(filepath.Join(m.shareDir, entry.Name()))
					os.Remove(filepath.Join(m.shareDir, info.Code+".zip"))
				}
			}
		}
	}
	return nil
}
