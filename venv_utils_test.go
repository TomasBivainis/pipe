package main

import (
	"os"
	"testing"
)

// Helper to create a temporary directory
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

func TestInstallUninstallDetectPackages(t *testing.T) {
	setupTempDirectory(t)

	// Create venv first
	err := createVirtualEnvironment()
	if err != nil {
		t.Skip("Could not create virtual environment (is python installed?):", err)
	}

	// Try installing a harmless package
	err = installPackages([]string{"wheel"})
	if err != nil {
		t.Fatalf("Could not install package (pip may not be available): %v", err)
	}

	packageExists, err := isPythonPackageInstalled("wheel")
	if err != nil {
		t.Fatalf("Could not check if package was installed: %v", err)
	}

	if !packageExists {
		t.Errorf("Package was not installed")
	}

	// Try uninstalling the package
	err = uninstallPackages([]string{"wheel"})
	if err != nil {
		t.Fatalf("Could not uninstall package (pip may not be available): %v", err)
	}

	packageExists, err = isPythonPackageInstalled("wheel")
	if err != nil {
		t.Fatalf("Could not check if package was installed: %v", err)
	}

	if packageExists {
		t.Errorf("Package was not uninstalled")
	}
}

// Try installing all the packages in the requirements file
func TestInstallPackagesFromRequirements(t *testing.T) {
	packages := []string{"requests", "pytest", "numpy"}

	setupTempRequirements(t, packages)

	// Create venv first
	err := createVirtualEnvironment()
	if err != nil {
		t.Skip("Could not create virtual environment (is python installed?):", err)
	}

	err = installPackagesFromRequirements()
	if err != nil {
		t.Errorf("Could not install packages from requirements: %v", err)
	}

	for _, p := range packages {
		installed, err := isPythonPackageInstalled(p)
		if err != nil {
			t.Fatalf("Could not check if package was installed: %v", err)
		}

		if !installed {
			t.Errorf("A package was not installed")
		}
	}
}

func TestGetVenvPythonPath(t *testing.T) {
	setupTempDirectory(t)

	// Create venv first
	err := createVirtualEnvironment()
	if err != nil {
		t.Skip("Could not create virtual environment (is python installed?):", err)
	}

	pipPath, err := getVenvPythonPath()
	if err != nil {
		t.Fatalf("getVenvPythonPath error: %v", err)
	}

	if pipPath == "" {
		t.Error("Expected a pip path, got empty string")
	}
}

func TestGetVenvPythonPathWithoutVirtualEnvironment(t *testing.T) {
	setupTempDirectory(t)

	_, err := getVenvPythonPath()
	if err == nil {
		t.Fatalf("Expected error when searching for venv python without virtual environment, but got none")
	}
}

func TestRunScript(t *testing.T) {
	setupTempDirectory(t)

	// Create venv first
	err := createVirtualEnvironment()
	if err != nil {
		t.Skip("Could not create virtual environment (is python installed?):", err)
	}

	// Create a simple Python script for testing
	scriptContent := `#!/usr/bin/env python3
import sys
print("Hello from Python script!")
print(f"Python version: {sys.version}")
sys.exit(0)
`
	
	scriptPath := "test_script.py"
	err = os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Test running the script
	err = runScript(scriptPath)
	if err != nil {
		t.Fatalf("runScript failed: %v", err)
	}
}

func TestRunScriptWithNonExistentScript(t *testing.T) {
	setupTempDirectory(t)

	// Create venv first
	err := createVirtualEnvironment()
	if err != nil {
		t.Skip("Could not create virtual environment (is python installed?):", err)
	}

	// Test running a non-existent script
	err = runScript("non_existent_script.py")
	if err == nil {
		t.Error("Expected error when running non-existent script, but got none")
	}
}

func TestRunScriptWithoutVirtualEnvironment(t *testing.T) {
	setupTempDirectory(t)

	// Test running a script without a virtual environment
	err := runScript("test_script.py")
	if err == nil {
		t.Error("Expected error when running script without virtual environment, but got none")
	}
}