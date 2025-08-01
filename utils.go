package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


func getGlobalPythonPath() (string, error) {
	candidates := []string{"python3", "python", "py"}

	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("python not found")
}

func getVenvPipPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	venvDirs := []string{"venv", ".venv", "env"}

	for _, venvDir := range venvDirs {
		// Adjust for OS
		venvPip := filepath.Join(cwd, venvDir, "bin", "pip")
		if _, err := os.Stat(venvPip); err == nil {
			return venvPip, nil
		}

		venvPipWin := filepath.Join(cwd, venvDir, "Scripts", "pip.exe")
		if _, err := os.Stat(venvPipWin); err == nil {
			return venvPipWin, nil
		}
	}

	

	return "", fmt.Errorf("pip not found in virtual environment")
}

func getFilePath(path string) (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	path = filepath.Join(cwd, path)

	if _, err := os.Stat(path) ; err == nil {
		return path, nil
	} else if os.IsNotExist(err) {
		return "", nil
	} else {
		return "", err
	}
}

func createRequirementsFile() error {
	targetFile := "requirements.txt"

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(cwd, targetFile))

	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func writePackagesToRequirementsFile(packages []string) error {
	requirementsFile, err := getFilePath("requirements.txt")
	if err != nil {
		return err
	}
	if requirementsFile == "" {
		return fmt.Errorf("requirements.txt not found")
	}

	content := strings.Join(packages, "\n")
	err = os.WriteFile(requirementsFile, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func getPackagesFromRequirements() ([]string, error) {
	requirementsFile, err := getFilePath("requirements.txt")
	if err != nil {
		return nil, err
	}
	if requirementsFile == "" {
		return nil, fmt.Errorf("requirements.txt not found")
	}

	data, err := os.ReadFile(requirementsFile)
	if err != nil {
			return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var packages []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			packages = append(packages, trimmed)
		}
	}
	return packages, nil
}

func detectVirtualEnvironment() (bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	venvDirs := []string{"venv", ".venv", "env"}
	pythonNames := []string{"bin/python", "Scripts/python.exe"}

	for _, venv := range venvDirs {
		venvPath := filepath.Join(cwd, venv)
		info, err := os.Stat(venvPath)
		if err == nil && info.IsDir() {
			// Check for python executable inside venv
			for _, py := range pythonNames {
				pyPath := filepath.Join(venvPath, py)
				if _, err := os.Stat(pyPath); err == nil {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func createVirtualEnvironment() error {
    pythonPath, err := getGlobalPythonPath()
    if err != nil {
        return err
    }
    cmd := exec.Command(pythonPath, "-m", "venv", ".venv")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func installPackages(packages []string) error {
	pipCommand, err := getVenvPipPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(pipCommand, (append([]string{"install"}, packages...))...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installPackagesFromRequirements() error {
	pipCommand, err := getVenvPipPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(pipCommand, "install", "-r", "requirements.txt")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func uninstallPackages(packages []string) error {
	pipCommand, err := getVenvPipPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(pipCommand, (append([]string{"uninstall"}, packages...))...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func removePackagesFromRequirementsFile(packages []string) error {
	currentPackages, err := getPackagesFromRequirements()
	if err != nil {
		return err
	}

	// Create a set of packages to remove (case-insensitive)
    toRemove := make(map[string]struct{})
    for _, pkg := range packages {
			toRemove[strings.ToLower(pkg)] = struct{}{}
    }

	// Filter out packages to remove
	var updatedPackages []string
	for _, pkg := range currentPackages {
		if _, found := toRemove[strings.ToLower(pkg)]; !found {
			updatedPackages = append(updatedPackages, pkg)
		}
	}

	writePackagesToRequirementsFile(updatedPackages)

	return nil
}

func addPackagesToRequirementsFile(packages []string) error {
	// Get already written packages
	currentPackages, err := getPackagesFromRequirements()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create a set for fast lookup
	existing := make(map[string]struct{})
	for _, pkg := range currentPackages {
		existing[strings.ToLower(pkg)] = struct{}{}
	}

	// Filter out duplicates
	var newPackages []string
	for _, pkg := range packages {
		if _, found := existing[strings.ToLower(pkg)]; !found {
			newPackages = append(newPackages, pkg)
			existing[strings.ToLower(pkg)] = struct{}{}
		}
	}

	if len(newPackages) == 0 {
		return nil // Nothing new to write
	}

	updatedPackages := append(currentPackages, newPackages...)

	writePackagesToRequirementsFile(updatedPackages)
	
	return nil
}