package main

import (
	"os"
	"testing"
)

func setupTempDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	
	t.Cleanup(func() { os.Chdir(oldDir) })
}

func TestGetGlobalPythonPath(t *testing.T) {
	path, err := getGlobalPythonPath()
	if err != nil {
		t.Skip("Python not found on system, skipping test")
	}
	
	if path == "" {
		t.Error("Expected a python path, got empty string")
	}
}

func TestCreateAndDetectVirtualEnvironment(t *testing.T) {
	setupTempDirectory(t)

	err := createVirtualEnvironment()
	if err != nil {
		t.Skip("Could not create virtual environment (is python installed?):", err)
	}

	found, err := detectVirtualEnvironment()
	if err != nil {
		t.Fatalf("detectVirtualEnvironment error: %v", err)
	}

	if !found {
		t.Error("Expected to detect a virtual environment, but did not")
	}
}

func TestGetVenvPipPath(t *testing.T) {
	setupTempDirectory(t)

	// Create venv first
	err := createVirtualEnvironment()
	if err != nil {
		t.Skip("Could not create virtual environment (is python installed?):", err)
	}

	pipPath, err := getVenvPipPath()
	if err != nil {
		t.Fatalf("getVenvPipPath error: %v", err)
	}

	if pipPath == "" {
		t.Error("Expected a pip path, got empty string")
	}
}

func TestInstallAndUninstallPackages(t *testing.T) {
	setupTempDirectory(t)

	// Create venv first
	err := createVirtualEnvironment()
	if err != nil {
		t.Skip("Could not create virtual environment (is python installed?):", err)
	}

	// Try installing a harmless package
	err = installPackages([]string{"wheel"})
	if err != nil {
		t.Errorf("Could not install package (pip may not be available): %v", err)
	}

	// Try uninstalling the package
	err = uninstallPackages([]string{"wheel"})
	if err != nil {
		t.Errorf("Could not uninstall package (pip may not be available): %v", err)
	}
}