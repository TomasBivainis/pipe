package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// returns the path to the file if found
// else returns empty string
func getFilePath(filename string) (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	path := filepath.Join(cwd, filename)

	if _, err := os.Stat(path) ; err == nil {
		return path, nil
	} else if os.IsNotExist(err) {
		return "", nil
	} else {
		return "", err
	}
}

// creates a requirements.txt file
func createRequirementsFile() error {
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

// writes the list of passed packages to the requirements.txt file
// overwrites the text inside the file
func writePackagesToRequirementsFile(packages []string) error {
	requirementsFile, err := getFilePath("requirements.txt")
	if err != nil {
		return err
	}
	if requirementsFile == "" {
		return fmt.Errorf("requirements.txt not found")
	}

	content := strings.Join(packages, "\n")
	err = os.WriteFile(requirementsFile, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

// returns a list of packages in the requirements.txt file
// as a list of strings
func getPackagesFromRequirements() ([]string, error) {
	requirementsFile, err := getFilePath("requirements.txt")
	if err != nil {
		return nil, err
	}
	if requirementsFile == "" {
		return nil, fmt.Errorf("requirements.txt not found")
	}

	data, err := os.ReadFile(requirementsFile)
	if err != nil {
			return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var packages []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			packages = append(packages, trimmed)
		}
	}
	return packages, nil
}

// removes the given list of packages from the requirements.txt file
func removePackagesFromRequirementsFile(packages []string) error {
	currentPackages, err := getPackagesFromRequirements()
	if err != nil {
		return err
	}

	// Create a set of packages to remove (case-insensitive)
    toRemove := make(map[string]struct{})
    for _, pkg := range packages {
			toRemove[strings.ToLower(pkg)] = struct{}{}
    }

	// Filter out packages to remove
	var updatedPackages []string
	for _, pkg := range currentPackages {
		if _, found := toRemove[strings.ToLower(pkg)]; !found {
			updatedPackages = append(updatedPackages, pkg)
		}
	}

	err = writePackagesToRequirementsFile(updatedPackages)
	if err != nil {
		return err
	}

	return nil
}

// adds the passed packages to the requirements.txt file
// duplicates do not get added again
func addPackagesToRequirementsFile(packages []string) error {
	// Get already written packages
	currentPackages, err := getPackagesFromRequirements()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create a set for fast lookup
	existing := make(map[string]struct{})
	for _, pkg := range currentPackages {
		existing[strings.ToLower(pkg)] = struct{}{}
	}

	// Filter out duplicates
	var newPackages []string
	for _, pkg := range packages {
		if _, found := existing[strings.ToLower(pkg)]; !found {
			newPackages = append(newPackages, pkg)
			existing[strings.ToLower(pkg)] = struct{}{}
		}
	}

	if len(newPackages) == 0 {
		return nil // Nothing new to write
	}

	updatedPackages := append(currentPackages, newPackages...)

	err = writePackagesToRequirementsFile(updatedPackages)
	if err != nil {
		return err
	}
	
	return nil
}