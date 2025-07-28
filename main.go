package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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

func detectRequirementsFile() (string, error) {
	targetFile := "requirements.txt"

	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	filepath := filepath.Join(cwd, targetFile)

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
	requirementsFile, err := detectRequirementsFile()

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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ami",
		Short: "ami is a simple package manager meant to improve pip",
		Long:  `A package manager CLI built in Go to replace pip.`,
	}

	// Add "init" command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing a python new project...")

			pythonCommand, err := detectPython()

			if err != nil {
				fmt.Println("Command error:", err)
				return
			}

			shellCmd := exec.Command(pythonCommand, "-m", "venv", ".venv")

			shellCmd.Stdout = os.Stdout
			shellCmd.Stderr = os.Stderr
			shellCmd.Stdin = os.Stdin

			existsRequirementsFile, err := detectRequirementsFile()

			if err != nil {
				fmt.Println("Command error:", err)
				return
			}

			if existsRequirementsFile == "" {
				err := createRequirementsFile()

				if err != nil {
					fmt.Println("Command error:", err)
					return;
				}
			}

			if err := shellCmd.Run(); err != nil {
				fmt.Println("Command failed:", err)
			}
		},
	})

	// Add adding to the requirements file
	rootCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install a python pip package",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("Pretending to install package:", args)

			pipCommand, err := detectPip()

			if err != nil {
				fmt.Println("Command failed:", err)
				return
			}

			if len(args) == 0 {
				//fmt.Println("install requirement package")

				shellCmd := exec.Command(pipCommand, "install", "-r", "requirements.txt")

				shellCmd.Stdout = os.Stdout
				shellCmd.Stderr = os.Stderr
				shellCmd.Stdin = os.Stdin

				if err := shellCmd.Run(); err != nil {
					fmt.Println("Command failed:", err)
				}
			} else {
				//fmt.Println("install named package")

				shellCmd := exec.Command(pipCommand, "install", strings.Join(args, " "))

				shellCmd.Stdout = os.Stdout
				shellCmd.Stderr = os.Stderr
				shellCmd.Stdin = os.Stdin

				if err := shellCmd.Run(); err != nil {
					fmt.Println("Command failed:", err)
					return
				}

				err := addPackagesToRequirementsFile(args)
				if err != nil {
					fmt.Println("Command error:", err)
				}
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}