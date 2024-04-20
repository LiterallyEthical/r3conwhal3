package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/LiterallyEthical/r3conwhal3/internal/mods"
	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/LiterallyEthical/r3conwhal3/pkg/logger"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cmds     = []string{"subfinder", "assetfinder", "amass", "httpx", "massdns", "puredns", "gotator", "gowitness", "ffuf"}
	myLogger = logger.GetLogger()
	//go:embed docs/*
	docFS          embed.FS
	specifiedFiles []string
)

func main() {

	defer utils.CleanUp()

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
	var enableAllMods, enablePassiveEnum, enableActiveEnum, enableWebOps bool

	pflag.StringVarP(&domain, "domain", "d", "", "Target domain to enumerate")
	pflag.StringVarP(&configDir, "config-dir", "c", "embedded", "Path to directory which config.env exists")
	pflag.StringVarP(&outDir, "out-dir", "o", "$HOME/user/r3conwhal3/results", "Directory to keep all output")
	pflag.BoolVarP(&enablePassiveEnum, "passive", "p", false, "Perform passsive subdomain enumeration process")
	pflag.BoolVarP(&enableActiveEnum, "active", "a", false, "Perform active recon processs (DNS bruteforce & DNS permutation)")
	pflag.BoolVarP(&enableAllMods, "all", "A", true, "Perform all passive & active recon process")
	pflag.BoolVarP(&enableWebOps, "webops", "w", false, "Perform web operations such as webscreenshoting, directory fuzzing etc.")
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

	// Set configs for PASSIVE_ENUM
	passiveEnumCFG := mods.PassiveEnum{
		Domain:            domain,
		FilePath:          passiveFilePath,
		OutDirPath:        outDirPath,
		EnableAssetfinder: config.EnableAssetfinder,
		EnableAmass:       config.EnableAmass,
		EnableSubkill3r:   config.EnableSubkill3r,
		Subfinder: mods.Subfinder{
			NumOfThreads: config.SubfinderNumOfThreads,
		},
		Amass: mods.Amass{
			Timeout: config.AmassTimeout,
		},
		Subkill3r: mods.Subkill3r{
			Wordlist:    config.Subkill3rWordlist,
			ServerAddr:  config.Subkill3rServerAddr,
			WorkerCount: config.Subkill3rWorkerCount,
		},
	}

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

	// Set configs for WEB_OPS
	webopsCFG := mods.WebOps{
		OutDirPath:      outDirPath,
		EnableGowitness: config.EnableGowitness,
		EnableFFUF:      config.EnableFFUF,
		Gowitness: mods.Gowitness{
			Timeout:               config.GowitnessTimeout,
			ResolutionX:           config.GowitnessResolutionX,
			ResolutionY:           config.GowitnessResolutionY,
			NumOfThreads:          config.GowitnessNumOfThreads,
			Fullpage:              config.GowitnessFullpage,
			ScreenshotFilter:      config.GowitnessScreenshotFilter,
			ScreenshotFilterCodes: config.GowitnessScreenshotFilterCodes,
		},
		FFUF: mods.FFUF{
			NumOfThreads:       config.FFUFNumOfThreads,
			Maxtime:            config.FFUFMaxtime,
			Rate:               config.FFUFRate,
			Timeout:            config.FFUFTimeout,
			Wordlist:           config.FFUFWordlist,
			MatchHTTPCode:      config.FFUFMatchHTTPCode,
			FilterResponseSize: config.FFUFFilterResponseSize,
			OutputFormat:       config.FFUFOutputFormat,
			Output:             config.FFUFOutput,
			SF:                 config.FFUFSF,
			SE:                 config.FFUFSE,
		},
	}

	// Run passive enumeration if enabled or no flags are provided (default behavior)
	if enablePassiveEnum || (!enableActiveEnum && !enablePassiveEnum) {
		if err := mods.InitSubdEnum(passiveEnumCFG); err != nil {
			log.Fatal(err)
		}
	}

	// Run active enumeration if enabled or no flags are provided (default behavior)
	if enableActiveEnum || (!enableActiveEnum && !enablePassiveEnum) {
		if err := mods.InitActiveSubdEnum(activeEnumCFG); err != nil {
			log.Fatal(err)
		}
	}

	if err := mods.InitFilterLiveDomains(outDirPath); err != nil {
		log.Fatal(err)
	}

	if enableWebOps || (!enableWebOps && !enableActiveEnum && !enablePassiveEnum) {
		if err := mods.InitWebOps(webopsCFG); err != nil {
			log.Fatal(err)
		}
	}
}
