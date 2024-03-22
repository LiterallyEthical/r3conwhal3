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
	cmds     = []string{"subfinder", "assetfinder", "amass", "httpx", "massdns", "puredns", "gotator"}
	myLogger = logger.GetLogger()
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

	// Check if the domain is provided or not
	if domain == "" {
		fmt.Println("Usage: go run main.go -d <domain> [-c <path-to-config-dir>] [-outDir <path-to-out-dir>]")
		return
	}

	// Check for installation of the required tools
	if err := utils.CheckInstallations(cmds); err != nil {
		log.Fatal(err)
	}

	// Create directory to keep all output
	outDirPath, err := utils.CreateDir(outDir, domain)
	if err != nil {
		log.Fatalf("Failed to create directory: %v, %v", outDir, err)
	}

	// Set filename and filepath for configs
	activeFileName := "active_enum_subdomains.txt"
	//activeFilePath := path.Join(outDirPath, activeFileName)
	passiveFileName := "passive_enum_subdomains.txt"
	passiveFilePath := path.Join(outDirPath, passiveFileName)
	specifiedFiles = append(specifiedFiles, passiveFileName, activeFileName)
	sublist := path.Join(outDirPath, "all_subdomains.txt")

	// Set configs for ACTIVE_ENUM
	activeEnumCFG := mods.ActiveEnum{
		PureDNS: mods.PureDNS{
			Domain:       domain,
			Wordlist:     config.PurednsWordlist,
			Resolvers:    config.PurednsResolvers,
			NumOfThreads: config.PurednsNumOfThreads,
		},
		Gotator: mods.Gotator{
			Sublist:      sublist,
			Permlist:     config.GotatorPermlist,
			Depth:        config.GotatorDepth,
			Numbers:      config.GotatorNumbers,
			NumOfThreads: config.GotatorNumOfThreads,
			Mindup:       config.GotatorMindup,
			Adv:          config.GotatorAdv,
			Md:           config.GotatorMd,
		},
		OutDirPath:     outDirPath,
		SpecifiedFiles: specifiedFiles,
	}

	if err := mods.InitSubdEnum(domain, passiveFilePath, outDirPath, wordlist, serverAddr, workerCount); err != nil {
		log.Fatal(err)
	}

	if err := mods.InitActiveSubdEnum(activeEnumCFG); err != nil {
		log.Fatal(err)
	}

	if err := mods.InitFilterLiveDomains(outDirPath); err != nil {
		log.Fatal(err)
	}
}
