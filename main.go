package main

import (
	"fmt"
	"os"
	"os/exec"

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

			shellCmd := exec.Command("python", "-m", "venv", ".venv")

			shellCmd.Stdout = os.Stdout
			shellCmd.Stderr = os.Stderr
			shellCmd.Stdin = os.Stdin

			if err := shellCmd.Run(); err != nil {
				fmt.Println("Command failed:", err)
			}
		},
	})

	// Add "install" command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install a package (not really)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Pretending to install package:", args)
		},
	})

	// Add "run" command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run a script",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running script:", args)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}