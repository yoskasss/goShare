package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_Create(t *testing.T) {
	tmpDir := t.TempDir()
	manager := &Manager{basePath: tmpDir}

	err := manager.Create("testproject", "Test summary")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	projectPath := filepath.Join(tmpDir, "testproject")
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Fatal("Project directory not created")
	}

	indexPath := filepath.Join(projectPath, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Fatal("index.html not created")
	}

	projectFile := filepath.Join(projectPath, ".goshare.json")
	if _, err := os.Stat(projectFile); os.IsNotExist(err) {
		t.Fatal(".goshare.json not created")
	}
}

func TestManager_CreateDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	manager := &Manager{basePath: tmpDir}

	err := manager.Create("testproject", "Test summary")
	if err != nil {
		t.Fatalf("First create failed: %v", err)
	}

	err = manager.Create("testproject", "Another summary")
	if err == nil {
		t.Fatal("Expected error for duplicate project")
	}
}

func TestManager_LoadProject(t *testing.T) {
	tmpDir := t.TempDir()
	manager := &Manager{basePath: tmpDir}

	err := manager.Create("testproject", "Test summary")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	project, err := manager.loadProject("testproject")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if project.Name != "testproject" {
		t.Errorf("Expected name 'testproject', got '%s'", project.Name)
	}
	if project.Summary != "Test summary" {
		t.Errorf("Expected summary 'Test summary', got '%s'", project.Summary)
	}
	if project.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestManager_List(t *testing.T) {
	tmpDir := t.TempDir()
	manager := &Manager{basePath: tmpDir}

	manager.Create("project1", "Summary 1")
	manager.Create("project2", "Summary 2")

	projects, err := manager.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
}

func TestManager_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	manager := &Manager{basePath: tmpDir}

	err := manager.Create("testproject", "Test summary")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = manager.Delete("testproject")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	projectPath := filepath.Join(tmpDir, "testproject")
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		t.Fatal("Project directory still exists after delete")
	}
}

func TestManager_GenerateFileTree(t *testing.T) {
	tmpDir := t.TempDir()
	manager := &Manager{basePath: tmpDir}

	os.MkdirAll(filepath.Join(tmpDir, "dir1", "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "dir1", "file2.txt"), []byte("content"), 0644)

	tree := manager.generateFileTree(tmpDir)

	if !contains(tree, "dir1") {
		t.Error("Tree should contain dir1")
	}
	if !contains(tree, "file1.txt") {
		t.Error("Tree should contain file1.txt")
	}
	if !contains(tree, "file2.txt") {
		t.Error("Tree should contain file2.txt")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func TestManager_Copy(t *testing.T) {
	tmpDir := t.TempDir()
	manager := &Manager{basePath: tmpDir}

	err := manager.Create("testproject", "Test summary")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	sourceDir := filepath.Join(tmpDir, "source")
	os.MkdirAll(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0644)

	err = manager.Copy("testproject", sourceDir)
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	destFile := filepath.Join(tmpDir, "testproject", "source", "test.txt")
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		t.Fatal("Copied file not found")
	}
}
