package main

import (
	"fmt"
	"os"

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
			}

			path, err := detectFile("requirements.txt")
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
		},
	})

	// Add adding to the requirements file
	rootCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install a python pip package",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("Pretending to install package:", args)

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
				err := installAllPackages()
				if err != nil {
					fmt.Println("Error while installing packages:", err)
					return
				}
			} else {
				err := installPackages(args)
				if err != nil {
					fmt.Println("Erro while installing packages:", err)
					return
				}

				err = writePackagesToRequirementsFile(args)
				if err != nil {
					fmt.Println("Erro while writing packages to requirements file:", err)
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