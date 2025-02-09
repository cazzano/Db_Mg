package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Store is the refactored main function from store.go
func Store() {
	// Expand the path to the config file
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "database", "data_base.json")

	// Read the config file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		// Initialize a new config if file doesn't exist
		configData = []byte("{\"location\":\"\", \"folders\":{}}")
	}

	// Parse the JSON config
	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return
	}

	// Ensure folders map is initialized
	if config.Folders == nil {
		config.Folders = make(map[string]Folder)
	}

	// Print the base location
	fmt.Printf("Base location from config: %s\n", config.Location)

	// Create a reader for input
	reader := bufio.NewReader(os.Stdin)

	// Prompt for folder name
	fmt.Print("Enter the name of the folder you want to create: ")
	folderName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	// Trim whitespace and newline
	folderName = strings.TrimSpace(folderName)

	// Validate folder name
	if folderName == "" {
		fmt.Println("Folder name cannot be empty.")
		return
	}

	// Prompt for folder description (optional)
	fmt.Print("Enter a description for the folder (optional): ")
	description, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading description: %v\n", err)
		return
	}
	description = strings.TrimSpace(description)

	// Construct full path
	fullPath := filepath.Join(config.Location, folderName)

	// Confirm folder creation
	fmt.Printf("Do you want to create folder '%s' at '%s'? (yes/no): ", folderName, fullPath)
	var response string
	fmt.Scanln(&response)

	// Check user response
	if response == "yes" || response == "y" {
		// Check if folder already exists
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			fmt.Printf("Folder '%s' already exists. Overwrite? (yes/no): ", folderName)
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Folder creation cancelled.")
				return
			}
		}

		// Create the directory
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}

		// Create folder entry in config
		folderEntry := Folder{
			Path:        fullPath,
			CreatedAt:   time.Now(),
			Description: description,
		}
		config.Folders[folderName] = folderEntry

		// Write updated config back to file
		updatedConfigData, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Printf("Error preparing config data: %v\n", err)
			return
		}

		err = os.WriteFile(configPath, updatedConfigData, 0644)
		if err != nil {
			fmt.Printf("Error updating config file: %v\n", err)
			return
		}

		fmt.Printf("Folder '%s' created successfully at: %s\n", folderName, fullPath)
		fmt.Println("Folder information saved to configuration.")

		// Ask the user if they want to drop files into the newly created folder
		fmt.Print("Do you want to drop files into this folder? (yes/no): ")
		fmt.Scanln(&response)
		if response == "yes" || response == "y" {
			// Call the drop logic
			DropFiles(fullPath)
		} else {
			fmt.Println("No files will be dropped into the folder.")
		}
	} else {
		fmt.Println("Folder creation cancelled.")
	}
}

// DropFiles calls the drop logic to copy or move files into the specified folder
func DropFiles(destinationFolder string) {
	// Set the current working directory to the destination folder
	originalDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	defer os.Chdir(originalDir) // Restore the original working directory

	// Change to the destination folder
	err = os.Chdir(destinationFolder)
	if err != nil {
		fmt.Printf("Error changing to directory '%s': %v\n", destinationFolder, err)
		return
	}

	// Call the drop logic
	fmt.Println("You can now drop files into the folder:")
	Drop()
}
