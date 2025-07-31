package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


func detectPython() (string, error) {
	candidates := []string{"python3", "python", "py"}

	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil;
		}
	}

	return "", fmt.Errorf("python not found")
}

func detectVevnPip() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Adjust for OS
	venvPip := filepath.Join(cwd, "venv", "bin", "pip")
	if _, err := os.Stat(venvPip); err == nil {
		return venvPip, nil
	}

	venvPipWin := filepath.Join(cwd, "venv", "Scripts", "pip.exe")
	if _, err := os.Stat(venvPipWin); err == nil {
		return venvPipWin, nil
	}

	return "", fmt.Errorf("pip not found in virtual environment")
}

func detectFile(path string) (string, error) {
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

func createRequirementsFile() (error) {
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

func writePackagesToRequirementsFile(packages []string) (error) {
	requirementsFile, err := detectFile("requirements.txt")

	if err != nil {
		return err
	}

	file, err := os.OpenFile(requirementsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	text := strings.Join(packages, "\n")

	if _, err := file.WriteString(text); err != nil {
		return err
	}

	return nil
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
    pythonPath, err := detectPython()
    if err != nil {
        return err
    }
    cmd := exec.Command(pythonPath, "-m", "venv", "venv")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func installPackages(packages []string) (error) {
	pipCommand, err := detectVevnPip()
	if err != nil {
		return err
	}

	cmd := exec.Command(pipCommand, (append([]string{"install"}, packages...))...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installAllPackages() (error) {
	pipCommand, err := detectVevnPip()
	if err != nil {
		return err
	}

	cmd := exec.Command(pipCommand, "install", "-r", "requirements.txt")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}