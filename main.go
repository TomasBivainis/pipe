package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "pvm",
		Short: "pvm helps with package management.",
		Long:  `pvm is a package manager CLI built to improve the usage of pip and python.`,
	}

	// init command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing a python new project...")

			virtualEnvironmentExists, err := detectVirtualEnvironment()
			if err != nil {
				fmt.Println("Error while detecting virtual environment:", err)
				return
			}

			if !virtualEnvironmentExists {
				err := createVirtualEnvironment()
				if err != nil {
					fmt.Println("Error while creating virtual environment:", err)
					return
				}

				fmt.Println("Created a new virtual environment.")
			}

			path, err := getFilePath("requirements.txt")
			if err != nil {
				fmt.Println("Error while detecting requirements file:", err)
				return
			}

			if path == "" {
				err := createRequirementsFile()
				if err != nil {
					fmt.Println("Error while creating requirements file:", err)
					return
				}
			}

			fmt.Println("Created a new requirements.txt file.")

			path, err = getFilePath(".gitignore")
			if err != nil {
				fmt.Println("Error while detecting gitignore file:", err)
				return
			}

			if path == "" {
				err := createGitignoreFile()
				if err != nil {
					fmt.Println("Error while creating gitignore file:", err)
					return
				}
			}

			fmt.Println("Created a new gitignore file.")
		},
	})

	// install command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install a python pip package",
		Run: func(cmd *cobra.Command, args []string) {
			virtualEnvironmentExists, err := detectVirtualEnvironment()
			if err != nil {
				fmt.Println("Error while detecting virtual environment:", err)
				return
			}

			if !virtualEnvironmentExists {
				fmt.Println("Virtual environment not initiated. Run \"pvm init\"")
				return
			}

			if len(args) == 0 {
				err := installPackagesFromRequirements()
				if err != nil {
					fmt.Println("Error while installing packages:", err)
					return
				}

				fmt.Println("All packages from the requirements file have been installed.")
			} else {
				err := installPackages(args)
				if err != nil {
					fmt.Println("Error while installing packages:", err)
					return
				}
				fmt.Println("The package(s) have been installed.")

				err = addPackagesToRequirementsFile(args)
				if err != nil {
					fmt.Println("Error while writing packages to requirements file:", err)
					return
				}
				fmt.Println("The package(s) have been written to the requirements file.")
			}
		},
	})

	// uninstall command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall a python pip package",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No packages entered to uninstall.")
				return
			}

			virtualEnvironmentExists, err := detectVirtualEnvironment()
			if err != nil {
				fmt.Println("Error while detecting virtual environment:", err)
				return
			}

			if !virtualEnvironmentExists {
				fmt.Println("Virtual environment not initiated. Run \"pvm init\"")
				return
			}

			err = uninstallPackages(args)
			if err != nil {
				fmt.Println("Error while uninstalling packages:", err)
				return
			}
			fmt.Println("The package(s) have been uninstalled.")
			
			err = removePackagesFromRequirementsFile(args)
			if err != nil {
				fmt.Println("Error while removing packages from requirements file:", err)
				return
			}
			fmt.Println("The package(s) have been removed from the requirements file.")
		},
	})

	// run command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Runs a specified python script in the virtual environment",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No scripts entered to run.")
				return
			}

			virtualEnvironmentExists, err := detectVirtualEnvironment()
			if err != nil {
				fmt.Println("Error while detecting virtual environment:", err)
				return
			}

			if !virtualEnvironmentExists {
				fmt.Println("Virtual environment not initiated. Run \"pvm init\"")
				return
			}

			scriptName := args[0]
			err = runScript(scriptName)
			if err != nil {
				fmt.Printf("Error while running script %s: %v\n", scriptName, err)
				return
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}