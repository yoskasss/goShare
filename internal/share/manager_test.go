package share

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"GoShare/internal/ziputil"
)

func TestManager_GenerateShareCode(t *testing.T) {
	tmpDir := t.TempDir()

	homeDir := tmpDir
	os.Setenv("HOME", homeDir)

	projectDir := filepath.Join(homeDir, "GoShareProjects", "testproject")
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "file.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(projectDir, ".goshare.json"), []byte(`{"name":"testproject"}`), 0644)

	manager := &Manager{
		shareDir:   filepath.Join(homeDir, "GoShareShares"),
		zipManager: ziputil.NewManager(),
	}

	code, err := manager.GenerateShareCode("testproject")
	if err != nil {
		t.Fatalf("GenerateShareCode failed: %v", err)
	}

	if len(code) != 8 {
		t.Errorf("Expected code length 8, got %d", len(code))
	}

	infoPath := filepath.Join(manager.shareDir, code+".json")
	if _, err := os.Stat(infoPath); os.IsNotExist(err) {
		t.Fatal("Share info file not created")
	}

	zipPath := filepath.Join(manager.shareDir, code+".zip")
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Fatal("Zip file not created")
	}
}

func TestManager_Download(t *testing.T) {
	tmpDir := t.TempDir()

	homeDir := tmpDir
	os.Setenv("HOME", homeDir)

	projectDir := filepath.Join(homeDir, "GoShareProjects", "testproject")
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "file.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(projectDir, ".goshare.json"), []byte(`{"name":"testproject"}`), 0644)

	manager := &Manager{
		shareDir:   filepath.Join(homeDir, "GoShareShares"),
		zipManager: ziputil.NewManager(),
	}

	code, err := manager.GenerateShareCode("testproject")
	if err != nil {
		t.Fatalf("GenerateShareCode failed: %v", err)
	}

	destDir := filepath.Join(homeDir, "GoShareProjects", "testproject")
	err = manager.Download(code, "")
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "file.txt")); os.IsNotExist(err) {
		t.Fatal("Downloaded file not found")
	}
}

func TestManager_DownloadExpired(t *testing.T) {
	tmpDir := t.TempDir()

	homeDir := tmpDir
	os.Setenv("HOME", homeDir)

	projectDir := filepath.Join(homeDir, "GoShareProjects", "testproject")
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "file.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(projectDir, ".goshare.json"), []byte(`{"name":"testproject"}`), 0644)

	manager := &Manager{
		shareDir:   filepath.Join(homeDir, "GoShareShares"),
		zipManager: ziputil.NewManager(),
	}

	code, err := manager.GenerateShareCode("testproject")
	if err != nil {
		t.Fatalf("GenerateShareCode failed: %v", err)
	}

	infoPath := filepath.Join(manager.shareDir, code+".json")
	data, _ := os.ReadFile(infoPath)
	var info ShareInfo
	json.Unmarshal(data, &info)
	info.ExpiresAt = time.Now().Add(-1 * time.Hour)
	data, _ = json.Marshal(info)
	os.WriteFile(infoPath, data, 0644)

	err = manager.Download(code, "")
	if err == nil {
		t.Fatal("Expected error for expired share code")
	}
}

func TestManager_CleanupExpired(t *testing.T) {
	tmpDir := t.TempDir()

	homeDir := tmpDir
	os.Setenv("HOME", homeDir)

	projectDir := filepath.Join(homeDir, "GoShareProjects", "testproject")
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "file.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(projectDir, ".goshare.json"), []byte(`{"name":"testproject"}`), 0644)

	manager := &Manager{
		shareDir:   filepath.Join(homeDir, "GoShareShares"),
		zipManager: ziputil.NewManager(),
	}

	code, err := manager.GenerateShareCode("testproject")
	if err != nil {
		t.Fatalf("GenerateShareCode failed: %v", err)
	}

	infoPath := filepath.Join(manager.shareDir, code+".json")
	data, _ := os.ReadFile(infoPath)
	var info ShareInfo
	json.Unmarshal(data, &info)
	info.ExpiresAt = time.Now().Add(-1 * time.Hour)
	data, _ = json.Marshal(info)
	os.WriteFile(infoPath, data, 0644)

	err = manager.CleanupExpired()
	if err != nil {
		t.Fatalf("CleanupExpired failed: %v", err)
	}

	if _, err := os.Stat(infoPath); !os.IsNotExist(err) {
		t.Fatal("Expired share info not cleaned up")
	}

	zipPath := filepath.Join(manager.shareDir, code+".zip")
	if _, err := os.Stat(zipPath); !os.IsNotExist(err) {
		t.Fatal("Expired zip not cleaned up")
	}
}
