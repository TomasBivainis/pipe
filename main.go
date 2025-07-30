package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

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
				fmt.Println("Python detection failed:", err)
				return
			}

			shellCmd := exec.Command(pythonCommand, "-m", "venv", ".venv")

			shellCmd.Stdout = os.Stdout
			shellCmd.Stderr = os.Stderr
			shellCmd.Stdin = os.Stdin

			requirementsFilePath, err := detectFile("requirements.txt")

			if err != nil {
				fmt.Println("File detection failed:", err)
				return
			}

			if requirementsFilePath == "" {
				err := createRequirementsFile()

				if err != nil {
					fmt.Println("Requirements file creation failed:", err)
					return;
				}
			}

			if err := shellCmd.Run(); err != nil {
				fmt.Println("Command failed:", err)
				return
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
				fmt.Println("Pip detection failed:", err)
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
					return
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
					fmt.Println("Editing requirements file failed:", err)
					return
				}
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}