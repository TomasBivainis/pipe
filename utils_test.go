package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper to create a temp requirements.txt file
func setupTempRequirements(t *testing.T, lines []string) string {
	tmpDir := t.TempDir()
	reqPath := filepath.Join(tmpDir, "requirements.txt")
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}
	if err := os.WriteFile(reqPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp requirements.txt: %v", err)
	}
	// Change working dir to tmpDir for the test
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(oldDir) })
	return reqPath
}

func TestWriteAndReadPackages(t *testing.T) {
	setupTempRequirements(t, []string{})
	pkgs := []string{"requests", "flask"}
	if err := writePackagesToRequirementsFile(pkgs); err != nil {
		t.Fatalf("writePackagesToRequirementsFile failed: %v", err)
	}
	readPkgs, err := getPackagesFromRequirements()
	if err != nil {
		t.Fatalf("getPackagesFromRequirements failed: %v", err)
	}
	if len(readPkgs) != 2 || readPkgs[0] != "requests" || readPkgs[1] != "flask" {
		t.Errorf("unexpected packages: %v", readPkgs)
	}
}

func TestAddPackagesToRequirementsFile(t *testing.T) {
	setupTempRequirements(t, []string{"requests"})
	err := addPackagesToRequirementsFile([]string{"flask", "requests"})
	if err != nil {
		t.Fatalf("addPackagesToRequirementsFile failed: %v", err)
	}
	pkgs, _ := getPackagesFromRequirements()
	if len(pkgs) != 2 {
		t.Errorf("expected 2 unique packages, got %v", pkgs)
	}
}

func TestRemovePackagesFromRequirementsFile(t *testing.T) {
	setupTempRequirements(t, []string{"requests", "flask", "pytest"})
	err := removePackagesFromRequirementsFile([]string{"flask"})
	if err != nil {
		t.Fatalf("removePackagesFromRequirementsFile failed: %v", err)
	}
	pkgs, _ := getPackagesFromRequirements()
	for _, p := range pkgs {
		if p == "flask" {
			t.Errorf("flask should have been removed")
		}
	}
}

func TestGetGlobalPythonPath(t *testing.T) {
	_, err := getGlobalPythonPath()
	if err != nil {
		t.Skip("Python not found on system, skipping test")
	}
}