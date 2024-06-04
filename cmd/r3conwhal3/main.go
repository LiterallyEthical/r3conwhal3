package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"

	"github.com/LiterallyEthical/r3conwhal3/internal/mods"
	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/LiterallyEthical/r3conwhal3/pkg/logger"
	"github.com/LiterallyEthical/r3conwhal3/web"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cmds     = []string{"subfinder", "assetfinder", "amass", "httpx", "massdns", "puredns", "gotator", "gowitness", "ffuf", "subzy"}
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

	// Define subcommands
	if len(os.Args) < 2 {
		fmt.Println("Usage: r3conwhal3 [run] [galery] options")
		os.Exit(1)
	}

	// Switch on the subcommand
	switch os.Args[1] {
	case "galery":
		handleGalery(os.Args[2:])
	case "run":
		handleRun(os.Args[2:])
	default:
		fmt.Println("expected 'galery' or 'run' subcommands")
		os.Exit(1)
	}
}

func handleGalery(args []string) {
	// Create a flag set for the galery subcommand
	galeryCmd := pflag.NewFlagSet("galery", pflag.ExitOnError)

	// Define flags for the galery subcommand
	var screenshotPath string
	galeryCmd.StringVarP(&screenshotPath, "path", "p", "", "Path to screenshots")

	// Parse the flags
	galeryCmd.Parse(args)

	// Ensure the path is provided
	if screenshotPath == "" {
		fmt.Println("Usage: r3conwhal3 galery -p <path-to-screenshots>")
		galeryCmd.PrintDefaults()
		return
	}

	// Init r3conwhal3 web galery
	myLogger.Info("Starting web server for gallery...")
	if err := web.StartServer(screenshotPath); err != nil {
		myLogger.Error("Web server error:", err)
	}
}

func handleRun(args []string) {
	// Define flags
	var domain, outDir, configDir string
	var enableAllMods, enablePassiveEnum, enableActiveEnum, enableWebOps, enableVulnScan bool

	runCmd := pflag.NewFlagSet("run", pflag.ExitOnError)
	runCmd.StringVarP(&domain, "domain", "d", "", "Target domain to enumerate")
	runCmd.StringVarP(&configDir, "config-dir", "c", "embedded", "Path to directory which config.env exists")
	runCmd.StringVarP(&outDir, "out-dir", "o", "$HOME/user/r3conwhal3/results", "Directory to keep all output")
	runCmd.BoolVarP(&enablePassiveEnum, "passive", "p", false, "Perform passive subdomain enumeration process")
	runCmd.BoolVarP(&enableActiveEnum, "active", "a", false, "Perform active recon process (DNS brute-force & DNS permutation)")
	runCmd.BoolVarP(&enableAllMods, "all", "A", true, "Perform all passive & active recon process")
	runCmd.BoolVarP(&enableWebOps, "webops", "w", false, "Perform web operations such as web screenshotting, directory fuzzing etc.")
	runCmd.BoolVarP(&enableVulnScan, "vulnscan", "v", false, "Perform vulnerability scanning")
	runCmd.Parse(args)

	// Check if the domain is provided or not
	if domain == "" {
		fmt.Println("Usage: r3conwhal3 run -d <domain> [-c <path-to-config-dir>] [-outDir <path-to-out-dir>]")
		runCmd.PrintDefaults()
		return
	}

	config, err := utils.LoadConfig(configDir, docFS)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Set the flag value from the config if not explicitly set via command line
	if !runCmd.Lookup("out-dir").Changed {
		// Get the value from config, if flag is not set
		outDir = viper.GetString("OUT_DIR")
	}

	// Binding variables from config.env to flags
	viper.BindPFlag("OUT_DIR", runCmd.Lookup("out-dir"))

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
		EnableWebGalery: config.EnableWebGalery,
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

	// Set configs for VULN_SCAN
	vulnScanCFG := mods.VulnScan{
		OutdirPath:  outDirPath,
		EnableSubzy: config.EnableSubzy,
		Subzy: mods.Subzy{
			Concurrency: config.SUBZYConcurrency,
			Timeout:     config.SUBZYTimeout,
			HideFails:   config.SUBZYHideFails,
			HTTPS:       config.SUBZYHTTPS,
			VerifySSL:   config.SUBZYVerifySSL,
			Vuln:        config.SUBZYVuln,
		},
	}

	// Channel to signal cleanup
	cleanupChan := make(chan struct{})

	// Channel to handle OS signals
	signalChan := make(chan os.Signal, 1)

	// Register for interrupt (CTRL+C) and termination signals
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Defer statement to ensure cleanup is called
	defer func() {
		myLogger.Info("Cleanup process is running...")
		utils.CleanUp()
		myLogger.Info("Cleanup complete, exiting...")
	}()

	// Ensure the cleanup channel is only closed once
	var once sync.Once
	closeCleanupChan := func() {
		once.Do(func() {
			close(cleanupChan)
		})
	}

	// Run the app
	go func() {
		if err := runApplication(enablePassiveEnum, enableActiveEnum, enableWebOps, enableVulnScan, passiveEnumCFG, activeEnumCFG, webopsCFG, vulnScanCFG, outDirPath, cleanupChan, closeCleanupChan); err != nil {
			myLogger.Error("Error while running r3conwhal3: %v", err)
			// Signal to cleanup
			closeCleanupChan()
		}
	}()

	select {
	case <-signalChan:
		fmt.Println()
		myLogger.Warning("Received interrupt signal, initiating cleanup...")
		closeCleanupChan()
	case <-cleanupChan:
		fmt.Println()
		myLogger.Warning("Received cleanup signal, initiating cleanup...")
	}
}

func runApplication(enablePassiveEnum, enableActiveEnum, enableWebOps, enableVulnScan bool, passiveEnumCFG mods.PassiveEnum, activeEnumCFG mods.ActiveEnum, webopsCFG mods.WebOps, vulnScanCFG mods.VulnScan, outDirPath string, cleanupChan chan struct{}, closeCleanupChan func()) error {
	defer closeCleanupChan()

	// Run passive enumeration if enabled or no flags are provided (default behavior)
	if enablePassiveEnum || (!enableActiveEnum && !enablePassiveEnum) {
		if err := mods.InitSubdEnum(passiveEnumCFG); err != nil {
			log.Fatal(err)
		}
	}

	// Run active enumeration if enabled or no flags are provided (default behavior)
	if enableActiveEnum || (!enableActiveEnum && !enablePassiveEnum) {
		if err := mods.InitActiveSubdEnum(activeEnumCFG); err != nil {
			myLogger.Error("Error in InitActiveSubdEnum:", err)
			return err
		}
	}

	if err := mods.InitFilterLiveDomains(outDirPath); err != nil {
		myLogger.Error("Error in InitFilterLiveDomains:", err)
		return err
	}

	if enableWebOps || (!enableWebOps && !enableActiveEnum && !enablePassiveEnum) {
		if err := mods.InitWebOps(webopsCFG); err != nil {
			myLogger.Error("Error in InitWebOps:", err)
			return err
		}
	}

	if enableVulnScan || (!enableVulnScan && !enableWebOps && !enableActiveEnum && !enablePassiveEnum) {
		if err := mods.InitVulnScan(vulnScanCFG); err != nil {
			myLogger.Error("Error in InitVulnScan:", err)
			return err
		}
	}

	if enableWebOps && webopsCFG.EnableGowitness && webopsCFG.EnableWebGalery {
		go func() {
			if err := mods.RunWebServer(outDirPath); err != nil {
				myLogger.Error("Web server error:", err)
				closeCleanupChan()
			}
		}()
		// Wait for the cleanup signal if the web server is running
		select {
		case <-cleanupChan:
			fmt.Println()
			myLogger.Warning("Cleanup signal received, stopping application tasks...")
			return nil
		}
	} else {
		// If web server is not running, return immediately after webopsCFG
		return nil
	}
}
