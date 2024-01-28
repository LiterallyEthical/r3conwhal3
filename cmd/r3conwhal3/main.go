package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/LiterallyEthical/r3conwhal3/internal/tools"
	"github.com/LiterallyEthical/r3conwhal3/pkg/logger"
)




var (
	cmds = []string{ "subfinder", "assetfinder", "amass", "httpx"}
	myLogger logger.Logger
)


func main() {

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Sprintf("Couldn't find user's home directory, %v", err))
		return
	}

	// Set the default value for dirName
	defaultDir := filepath.Join(homeDir, "r3conwhal3", "results")


	// Define flags
	var domain, fileName, dirName string
	wordlist := filepath.Join("../", "docs", "subdomains-1000.txt")
	workerCount := 1000
	serverAddr := "8.8.8.8:53"
	flag.StringVar(&domain, "domain",  "", "Target domain to enumerate")
	flag.StringVar(&fileName, "file-name", "subdomains.txt", "File to write subdomains")
	flag.StringVar(&dirName, "dir-name", defaultDir, "Directory to keep all output")
	flag.Parse()

	
	tools.Banner()


	// Check if the domain is provided or not
	if domain == "" {
		fmt.Println("Usage: go run main.go -domain <domain> [-fileName <fileName>] [-dirName <dirName>]")
		return
	}

	// Check for installation of the required tools
	if err := tools.CheckInstallations(cmds); err != nil {
		log.Fatal(err)
	}

	// Create directory to keep all output
	dirPath, err := tools.CreateDir(dirName, domain)
	if err != nil {
		myLogger.Error("Failed to create directory: %v, %v", dirName, err)
	}

	// Join full path
	filePath := path.Join(dirPath, fileName)
	
	if err := tools.InitSubdEnum(domain, filePath, dirPath, wordlist, serverAddr, workerCount); err != nil {
		log.Fatal(err)
	}
	
	
}