package main

import (
	"embed"
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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cmds = []string{ "subfinder", "assetfinder", "amass", "httpx"}
	myLogger logger.Logger
	//go:embed docs/banner.txt docs/subdomains-1000.txt
	content embed.FS
)




func main() {


	// Accessing files from the embedded docs directory
	data, err := fs.ReadFile(content, "docs/banner.txt")
	if err != nil {
		log.Panic("Error reading banner.txt:", err)
		os.Exit(1)
	}

	// Print the banner
	fmt.Println(color.CyanString(string(data)))
	
	// Define flags
	var domain, fileName, outDir, configDir string


	pflag.StringVarP(&domain, "domain", "d", "", "Target domain to enumerate")
    pflag.StringVarP(&configDir, "config-dir", "c", "", "Path to directory which config file(config.env) exists")
    pflag.StringVarP(&outDir, "out-dir", "o", "$HOME/user/r3conwhal3/results", "Directory to keep all output")
    pflag.StringVarP(&fileName, "file-name", "f", "subdomains.txt", "File to write subdomains")
	pflag.Parse()


	config, err := utils.LoadConfig(configDir)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

    // Set the flag value from the config if not explicitly set via command line
    if !pflag.Lookup("out-dir").Changed {
		// Get the value from config, if flag is not set
        outDir = viper.GetString("OUT_DIR") 
    }

	if !pflag.Lookup("file-name").Changed {
        fileName = viper.GetString("FILE_NAME") 
    }

	// Binding variables from config.env to flags	
	viper.BindPFlag("OUT_DIR", pflag.Lookup("out-dir"))
	viper.BindPFlag("FILE_NAME", pflag.Lookup("file-name"))


	// Define variables for subkill3r
	workerCount, err := strconv.Atoi(config.Subkill3rWorkerCount)
	if err != nil {
		panic(err)
		//myLogger.Error("error running subkill3r: workerCount is type string instead of int" , err)
	}
	serverAddr := config.Subkill3rServerAddr
	wordlist := config.Subkill3rWordlist


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