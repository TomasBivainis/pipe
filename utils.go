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

func detectPip() (string, error) {
	candidates := []string{"pip", "pip3"}

	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil;
		}
	}

	pythonPath, err := detectPython()

	if err == nil {
		return pythonPath + " -m pip", nil
	}

	return "", fmt.Errorf("pip not found")
}

func detectFile(filename string) (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	filepath := filepath.Join(cwd, filename)

	if _, err := os.Stat(filepath) ; err == nil {
		return filepath, nil
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

func addPackagesToRequirementsFile(packages []string) (error) {
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