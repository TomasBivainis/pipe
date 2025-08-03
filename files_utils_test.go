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

	packages := []string{"requests", "flask"}
	if err := writePackagesToRequirementsFile(packages); err != nil {
		t.Fatalf("writePackagesToRequirementsFile failed: %v", err)
	}

	readPackages, err := getPackagesFromRequirements()
	if err != nil {
		t.Fatalf("getPackagesFromRequirements failed: %v", err)
	}

	if len(readPackages) != 2 || readPackages[0] != "requests" || readPackages[1] != "flask" {
		t.Errorf("unexpected packages: %v", readPackages)
	}

	packages = append(packages, "numpy")

	if err := writePackagesToRequirementsFile(packages); err != nil {
		t.Fatalf("writePackagesToRequirementsFile failed: %v", err)
	}

	readPackages, err = getPackagesFromRequirements()
	if err != nil {
		t.Fatalf("getPackagesFromRequirements failed: %v", err)
	}

	if len(readPackages) != 3 || readPackages[0] != "requests" || readPackages[1] != "flask" || readPackages[2] != "numpy" {
		t.Errorf("unexpected packages: %v", readPackages)
	}
}

func TestAddPackagesToRequirementsFile(t *testing.T) {
	setupTempRequirements(t, []string{"requests"})

	err := addPackagesToRequirementsFile([]string{"flask", "requests", "numpy"})
	if err != nil {
		t.Fatalf("addPackagesToRequirementsFile failed: %v", err)
	}

	packages, err := getPackagesFromRequirements()
	if err != nil {
		t.Fatalf("getPackagesFromRequirements failed: %v", err)
	}

	if len(packages) != 3 || packages[0] != "requests" || packages[1] != "flask" || packages[2] != "numpy" {
		t.Errorf("expected 3 unique packages, got %v", packages)
	}
}

func TestRemovePackagesFromRequirementsFile(t *testing.T) {
	setupTempRequirements(t, []string{"requests", "flask", "pytest", "numpy"})

	err := removePackagesFromRequirementsFile([]string{"flask", "numpy", "flask"})
	if err != nil {
		t.Fatalf("removePackagesFromRequirementsFile failed: %v", err)
	}

	packages, err := getPackagesFromRequirements()
	if err != nil {
		t.Fatalf("getPackagesFromRequirements failed: %v", err)
	}

	for _, p := range packages {
		if p == "flask" || p == "numpy" {
			t.Errorf("package was not removed")
		}
	}

	if len(packages) != 2 || packages[0] != "requests" || packages[1] != "pytest" {
		t.Errorf("unexpected packages: %v", packages)
	}
}

func TestGetFilePathWhenFileExists(t *testing.T) {
	reqPath := setupTempRequirements(t, []string{})
	// add /private because for some reason path has it
	reqPath = filepath.Join("/private", reqPath)

	path, err := getFilePath("requirements.txt")
	if err != nil {
		t.Fatalf("getFilePath failed: %v", err)
	}

	if path != reqPath {
		t.Errorf("file path is incorrect: expected %s, received %s", reqPath, path)
	}
}

func TestGetFilePathWhenFileDoesNotExists(t *testing.T) {
	tmpDir := t.TempDir()

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(oldDir) })

	path, err := getFilePath("requirements.txt")
	if err != nil {
		t.Fatalf("getFilePath failed: %v", err)
	}

	if path != "" {
		t.Errorf("unexpected path was found")
	}
}

func TestCreateRequirementsFile(t *testing.T) {
	tmpDir := t.TempDir()

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(oldDir) })

	expectedPath := filepath.Join("/private", tmpDir, "requirements.txt")

	err := createRequirementsFile()
	if err != nil {
		t.Fatalf("createRequirementsFile failed: %v", err)
	}

	path, err := getFilePath("requirements.txt")
	if err != nil {
		t.Fatalf("getFilePath failed: %v", err)
	}

	if path == "" {
		t.Errorf("path not found")
	}

	// add /private to tmpDir because for some reason path has it
	if path != expectedPath {
		t.Errorf("path is incorrect: excpected %s, received %s", expectedPath, path)
	}
}

/*
func TestGetGlobalPythonPath(t *testing.T) {
	_, err := getGlobalPythonPath()
	if err != nil {
		t.Skip("Python not found on system, skipping test")
	}
}
*/