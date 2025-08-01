package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ami",
		Short: "ami helps with package management.",
		Long:  `A package manager CLI built to improve the usage of pip and python.`,
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

				fmt.Println("Created a new requirements.txt file.")
			}
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
				fmt.Println("Virtual environment not initiated. Run \"ami init\"")
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

				fmt.Println("Packages written to the requirements file.")
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
				fmt.Println("Virtual environment not initiated. Run \"ami init\"")
				return
			}

			err = uninstallPackages(args)
			if err != nil {
				fmt.Println("Error while uninstalling packages:", err)
				return
			}

			err = removePackagesFromRequirementsFile(args)
			if err != nil {
				fmt.Println("Error while removing packages from requirements file:", err)
				return
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}