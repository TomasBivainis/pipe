package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	cmd := exec.Command(pipCommand, (append([]string{"uninstall", "-y"}, packages...))...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

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

func activateVirtualEnvironment() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	venvPath := ""
	venvDirs := []string{"venv", ".venv", "env"}

	for _, venv := range venvDirs {
		venvPath = filepath.Join(cwd, venv)

		info, err := os.Stat(venvPath)
		if err == nil && info.IsDir() {
			break
		}
	}
	if venvPath == "" {
		return fmt.Errorf("virtual environment was not found")
	}

	// Check if activation script exists
	var activateScript string
	switch runtime.GOOS {
	case "windows":
		activateScript = filepath.Join(venvPath, "Scripts", "activate.bat")
		if _, err := os.Stat(activateScript); err != nil {
			return fmt.Errorf("activation script not found at %s", activateScript)
		}
	default:
		activateScript = filepath.Join(venvPath, "bin", "activate")
		if _, err := os.Stat(activateScript); err != nil {
			return fmt.Errorf("activation script not found at %s", activateScript)
		}
	}

	// Get the user's shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		// Fallback shells
		switch runtime.GOOS {
		case "windows":
			shell = "cmd.exe"
		default:
			shell = "/bin/bash"
		}
	}

	// Print activation instructions
	fmt.Printf("To activate the virtual environment, run:\n")
	switch runtime.GOOS {
	case "windows":
		fmt.Printf("  %s\n", activateScript)
	default:
		fmt.Printf("  source %s\n", activateScript)
	}
	fmt.Printf("\nOr to activate and run a command:\n")
	switch runtime.GOOS {
	case "windows":
		fmt.Printf("  cmd /c \"%s && your_command_here\"\n", activateScript)
	default:
		fmt.Printf("  source %s && your_command_here\n", activateScript)
	}

	return nil
}