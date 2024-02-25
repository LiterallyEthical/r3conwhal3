package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/LiterallyEthical/r3conwhal3/internal/mods"
	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/LiterallyEthical/r3conwhal3/pkg/logger"
	"github.com/fatih/color"
)

var (
	cmds = []string{ "subfinder", "assetfinder", "amass", "httpx"}
	myLogger logger.Logger
	//go:embed docs/banner.txt docs/subdomains-1000.txt
	content embed.FS
)




func main() {

	config, err := utils.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Accessing files from the embedded docs directory
	data, err := fs.ReadFile(content, "docs/banner.txt")
	if err != nil {
		log.Panic("Error reading banner.txt:", err)
		os.Exit(1)
	}

	// Print the banner
	fmt.Println(color.CyanString(string(data)))
	
	// Define flags
	var domain, fileName, outDir string

	// Define variables for subkill3r
	workerCount, err := strconv.Atoi(config.Subkill3rWorkerCount)
	if err != nil {
		panic(err)
		//myLogger.Error("error running subkill3r: workerCount is type string instead of int" , err)
	}
	serverAddr := config.Subkill3rServerAddr
	wordlist := config.Subkill3rWordlist

	flag.StringVar(&domain, "domain",  "", "Target domain to enumerate")
	flag.StringVar(&outDir, "out-dir", config.OutDir, "Directory to keep all output")
	flag.StringVar(&fileName, "file-name", config.FileName, "File to write subdomains")
	flag.Parse()


	// Check if the domain is provided or not
	if domain == "" {
		fmt.Println("Usage: go run main.go -domain <domain> [-fileName <fileName>] [-outDir <outDir>]")
		return
	}

	// Check for installation of the required tools
	if err := utils.CheckInstallations(cmds); err != nil {
		log.Fatal(err)
	}

	// Create directory to keep all output
	dirPath, err := utils.CreateDir(outDir, domain)
	if err != nil {
		myLogger.Error("Failed to create directory: %v, %v", outDir, err)
	}

	// Join full path
	filePath := path.Join(dirPath, fileName)
	
	if err := mods.InitSubdEnum(domain, filePath, dirPath, wordlist, serverAddr, workerCount); err != nil {
		log.Fatal(err)
	}

}