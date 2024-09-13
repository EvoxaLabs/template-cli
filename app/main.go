package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	useTypeScript bool
)

var rootCmd = &cobra.Command{
	Use:   "node-template-cli",
	Short: "A CLI for generating full-stack Node.js project templates",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a full-stack project template",
	Run: func(cmd *cobra.Command, args []string) {

		checkDependencies()
		selectProjectExt()
		selectProjectType()

	},
}

func init() {
	generateCmd.Flags().BoolVarP(&useTypeScript, "typescript", "t", false, "Generate project with TypeScript")
	rootCmd.AddCommand(generateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Check if Node.js is installed
func checkDependencies() {
	fmt.Println("Checking for Node.js...")

	if _, err := exec.LookPath("node"); err != nil {
		fmt.Println("Node.js is not installed.")
		fmt.Println("Please install Node.js from https://nodejs.org/en/download/")
		os.Exit(1)
	} else {
		fmt.Println("Node.js is installed.")
	}
}

// Select frontend, backend, or full-stack
func selectProjectType() {
	fmt.Println("Select project type:")
	fmt.Println("1. Frontend only")
	fmt.Println("2. Backend only")
	fmt.Println("3. Full-stack (Frontend + Backend)")

	var projectType int
	fmt.Scan(&projectType)

	switch projectType {
	case 1:
		selectFrontendFramework()
	case 2:
		selectBackendFramework()
	case 3:
		selectFrontendFramework()
		selectBackendFramework()
	default:
		fmt.Println("Invalid option. Please select a valid number.")
		selectProjectType()
	}
}

func selectProjectExt() {
	fmt.Println("Select project extension:")
	fmt.Println("1. TypeScript")
	fmt.Println("2. JavaScript")

	var projectExt int
	fmt.Scan(&projectExt)

	switch projectExt {
	case 1:
		useTypeScript = true
	case 2:
		useTypeScript = false
	default:
		fmt.Println("Invalid option. Please select a valid number.")
		selectProjectExt()
	}
}

// Select frontend framework (React, Next.js)
func selectFrontendFramework() {
	fmt.Println("Select a frontend framework:")
	fmt.Println("1. React")
	fmt.Println("2. Next.js")

	var framework int
	fmt.Scan(&framework)

	projectName := "frontend"
	switch framework {
	case 1:
		generateFrontendStructure(projectName, "react")
	case 2:
		generateFrontendStructure(projectName, "next")
	default:
		fmt.Println("Invalid option. Please select a valid number.")
		selectFrontendFramework()
	}
}

// Select backend framework (Node.js, Express)
func selectBackendFramework() {
	fmt.Println("Select a backend framework:")
	fmt.Println("1. Node.js with Express")

	var framework int
	fmt.Scan(&framework)

	projectName := "backend"
	switch framework {
	case 1:
		generateBackendStructure(projectName)
	default:
		fmt.Println("Invalid option. Please select a valid number.")
		selectBackendFramework()
	}
}

// Generate Frontend Structure (React or Next.js)
func isDirectoryEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		log.Fatalf("Error opening directory: %v", err)
	}
	defer f.Close()

	// Check if directory is empty
	_, err = f.Readdirnames(1)

	return err == io.EOF // Directory is empty

}

func generateFrontendStructure(basePath, framework string) {
	// Check if directory is empty
	if !isDirectoryEmpty(basePath) {
		fmt.Printf("Directory %s contains files that could conflict. Please use an empty directory.\n", basePath)
		os.Exit(1)
	}

	// Let npx handle app creation
	switch framework {
	case "react":
		cmd := exec.Command("npx", "create-react-app", basePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Error creating React app: %v", err)
		}

	case "next":
		cmd := exec.Command("npx", "create-next-app", basePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Error creating Next.js app: %v", err)
		}
	}

	fmt.Printf("Frontend project created at: %s\n", basePath)
}

// Generate Backend Structure (Node.js, Express)
func generateBackendStructure(basePath string) {
	dirs := []string{
		filepath.Join(basePath, "src", "controllers"),
		filepath.Join(basePath, "src", "routes"),
		filepath.Join(basePath, "src", "models"),
	}

	var ext string
	if useTypeScript {
		ext = ".ts"
	} else {
		ext = ".js"
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory %s: %v", dir, err)
		}
		fmt.Printf("Created directory: %s\n", dir)
	}

	files := map[string]string{
		filepath.Join(basePath, "src", "index"+ext): sampleExpressApp(),
		filepath.Join(basePath, "package.json"):     expressPackageJSON(),
	}

	createFiles(files)
}

// Create files with given content
func createFiles(files map[string]string) {
	for path, content := range files {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			log.Fatalf("Error creating file %s: %v", path, err)
		}
		fmt.Printf("Created file: %s\n", path)
	}
}

// Sample Express app file
func sampleExpressApp() string {
	return `const express = require('express');
const app = express();

app.get('/', (req, res) => {
  res.send('Hello from Express!');
});

app.listen(3000, () => {
  console.log('Server is running on port 3000');
});`
}

// React package.json

// Express package.json
func expressPackageJSON() string {
	return `{
  "name": "express-app",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.17.1"
  }
}`
}
