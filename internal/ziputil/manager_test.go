package ziputil

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestManager_CreateZip(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager()

	sourceDir := filepath.Join(tmpDir, "source")
	os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "subdir", "file2.txt"), []byte("content2"), 0644)

	zipPath := filepath.Join(tmpDir, "test.zip")
	err := manager.CreateZip(sourceDir, zipPath)
	if err != nil {
		t.Fatalf("CreateZip failed: %v", err)
	}

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Fatal("Zip file not created")
	}

	files, err := manager.ListZip(zipPath)
	if err != nil {
		t.Fatalf("ListZip failed: %v", err)
	}

	expectedFiles := []string{"file1.txt", "subdir/file2.txt"}
	for _, expected := range expectedFiles {
		found := false
		for _, f := range files {
			if f == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in zip", expected)
		}
	}
}

func TestManager_ExtractZip(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager()

	sourceDir := filepath.Join(tmpDir, "source")
	os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "subdir", "file2.txt"), []byte("content2"), 0644)

	zipPath := filepath.Join(tmpDir, "test.zip")
	err := manager.CreateZip(sourceDir, zipPath)
	if err != nil {
		t.Fatalf("CreateZip failed: %v", err)
	}

	destDir := filepath.Join(tmpDir, "dest")
	err = manager.ExtractZip(zipPath, destDir)
	if err != nil {
		t.Fatalf("ExtractZip failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "file1.txt")); os.IsNotExist(err) {
		t.Fatal("file1.txt not extracted")
	}

	if _, err := os.Stat(filepath.Join(destDir, "subdir", "file2.txt")); os.IsNotExist(err) {
		t.Fatal("file2.txt not extracted")
	}

	content, err := os.ReadFile(filepath.Join(destDir, "file1.txt"))
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}
	if string(content) != "content1" {
		t.Errorf("Expected 'content1', got '%s'", string(content))
	}
}

func TestManager_ExtractZip_Security(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager()

	zipPath := filepath.Join(tmpDir, "malicious.zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)
	writer, err := zipWriter.Create("../../../etc/passwd")
	if err != nil {
		t.Fatalf("Failed to create zip entry: %v", err)
	}
	writer.Write([]byte("malicious"))
	zipWriter.Close()
	zipFile.Close()

	destDir := filepath.Join(tmpDir, "dest")
	err = manager.ExtractZip(zipPath, destDir)
	if err == nil {
		t.Fatal("Expected error for path traversal attempt")
	}
}
