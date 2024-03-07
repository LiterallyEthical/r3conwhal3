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
	cmds     = []string{"subfinder", "assetfinder", "amass", "httpx", "puredns"}
	myLogger logger.Logger
	//go:embed docs/*
	docFS          embed.FS
	specifiedFiles []string
)

func main() {

	// Accessing files from the embedded docs directory
	data, err := fs.ReadFile(docFS, "docs/banner.txt")
	if err != nil {
		log.Panic("Error reading banner.txt:", err)
		os.Exit(1)
	}

	// Print the banner
	fmt.Println(color.CyanString(string(data)))

	// Define flags
	var domain, outDir, configDir string

	pflag.StringVarP(&domain, "domain", "d", "", "Target domain to enumerate")
	pflag.StringVarP(&configDir, "config-dir", "c", "embedded", "Path to directory which config.env exists")
	pflag.StringVarP(&outDir, "out-dir", "o", "$HOME/user/r3conwhal3/results", "Directory to keep all output")
	pflag.Parse()

	config, err := utils.LoadConfig(configDir, docFS)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Set the flag value from the config if not explicitly set via command line
	if !pflag.Lookup("out-dir").Changed {
		// Get the value from config, if flag is not set
		outDir = viper.GetString("OUT_DIR")
	}

	// Binding variables from config.env to flags
	viper.BindPFlag("OUT_DIR", pflag.Lookup("out-dir"))

	// Define variables for subkill3r
	workerCount, err := strconv.Atoi(config.Subkill3rWorkerCount)
	if err != nil {
		myLogger.Error("error running subkill3r: workerCount is type string instead of int", err)
	}
	serverAddr := config.Subkill3rServerAddr
	wordlist := config.Subkill3rWordlist

	// Define variables for puredns
	purednsWordlist := config.PurednsWordlist
	purednsResolvers := config.PurednsResolvers
	purednsNumOfThreads := config.PurednsNumOfThreads

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
	outDirPath, err := utils.CreateDir(outDir, domain)
	if err != nil {
		myLogger.Error("Failed to create directory: %v, %v", outDir, err)
	}
	// Join full path for activeSubdEnum
	activeFileName := "active_enum_subdomains.txt"
	activeFilePath := path.Join(outDirPath, activeFileName)
	passiveFileName := "passive_enum_subdomains.txt"
	passiveFilePath := path.Join(outDirPath, passiveFileName)
	specifiedFiles = append(specifiedFiles, passiveFileName, activeFileName)

	if err := mods.InitSubdEnum(domain, passiveFilePath, outDirPath, wordlist, serverAddr, workerCount); err != nil {
		log.Fatal(err)
	}

	if err := mods.InitActiveSubdEnum(outDirPath, domain, activeFilePath, purednsWordlist, purednsResolvers, purednsNumOfThreads); err != nil {
		log.Fatal(err)
	}

	if err := mods.InitFilterLiveDomains(outDirPath, specifiedFiles); err != nil {
		log.Fatal(err)
	}
}
