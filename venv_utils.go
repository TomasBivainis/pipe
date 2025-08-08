package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// returns the path to the global python application
func getGlobalPythonPath() (string, error) {
	candidates := []string{"python3", "python", "py"}

	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("python not found")
}

func getVenvPythonPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	venvDirs := []string{"venv", ".venv", "env"}

	for _, venvDir := range venvDirs {
		// Adjust for OS
		venvPython := filepath.Join(cwd, venvDir, "bin", "python")
		if _, err := os.Stat(venvPython); err == nil {
			return venvPython, nil
		}

		venvPythonWin := filepath.Join(cwd, venvDir, "Scripts", "python.exe")
		if _, err := os.Stat(venvPythonWin); err == nil {
			return venvPythonWin, nil
		}
	}

	return "", fmt.Errorf("python not found in virtual environment")
}

// returns the path to the virtual environments pip
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

// returns true if a virtual environment was 
// initiated in the current working directory
// otherwise, returns false
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

// creates a virtual environment
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

// installs the passed list of packages and writes new packages
// to the requirements.txt file
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

// installs all of the packages named in the requirements.txt file
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

// uninstalls the given list of packages and removes them
// from the requirements.txt file
func uninstallPackages(packages []string) error {
	pipCommand, err := getVenvPipPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(pipCommand, (append([]string{"uninstall", "-y"}, packages...))...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// returns true if the passed python package is installed
// false otherwise
func isPythonPackageInstalled(pkg string) (bool, error) {
	pipPath, err := getVenvPipPath()
	if err != nil {
		return false, err
	}

	cmd := exec.Command(pipPath, "show", pkg)
	err = cmd.Run()
	if err == nil {
		return true, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return false, nil // Not installed
	}
	return false, err // Some other error
}

func runScript(scriptName string) error {
	pythonPath, err := getVenvPythonPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(pythonPath, scriptName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}