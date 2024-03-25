package mods

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/LiterallyEthical/r3conwhal3/pkg/logger"
	"github.com/LiterallyEthical/r3conwhal3/pkg/subkill3r"

	"github.com/fatih/color"
)

var (
	subCount int
	myLogger = logger.GetLogger()
)

type PassiveEnum struct {
	Domain     string
	FilePath   string
	OutDirPath string
	Subfinder  Subfinder
	Amass      Amass
	Subkill3r  Subkill3r
}

type Subfinder struct {
	NumOfThreads int
}

type Amass struct {
	Timeout int
}

type Subkill3r struct {
	Wordlist    string
	ServerAddr  string
	WorkerCount int
}

func RunSubfinder(domain, filePath string, numOfThreads int) error {
	// fmt.Printf("\n[+]Starting subfinder\n")
	myLogger.Info("Running subfinder")

	// printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "subfinder")

	// Show progress
	utils.ShowProgress()

	// Run subfinder
	_, err := utils.RunCommand("subfinder", "-d", domain, "-o", filePath, "-t", numOfThreads)
	if err != nil {
		return err
	}

	// Count enumareted subdomains
	subCount, err = utils.CountLines(filePath)
	if err != nil {
		// log.Printf("Error measuring enumerated subdomains: %v ", err)
		myLogger.Warning("Failed to measure number of gathered subdomains: %v", err)
	}
	// fmt.Printf("\r[+]%v subdomains gathered", countedLines)
	myLogger.Info("%v subdomain found!", subCount)

	// Log process completion and elapsed time
	// fmt.Printf("\n[+]Subfinder executed successfully")
	myLogger.Info("Subfinder executed successfully")
	return nil
}

func RunAssetfinder(domain, filePath string) error {
	myLogger.Info("Running assetfinder")

	// printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "assetfinder")

	// Show progress
	utils.ShowProgress()

	// Run assetfinder
	output, err := utils.RunCommand("assetfinder", "-subs-only", domain)
	if err != nil {
		return err
	}

	// Write output to specified file
	err = utils.AppendToFile(filePath, output)
	if err != nil {
		//log.Printf("Error appending to file %s: %v", filePath, err)
		myLogger.Warning("Error appending to file %s: %v", filePath, err)
	}

	// Count enumareted subdomains

	oldSubCount := subCount
	subCount, err = utils.CountLines(filePath)
	if err != nil {
		//log.Printf("Error measuring enumerated subdomains: %v ", err)
		myLogger.Warning("Error measuring enumerated subdomains: %v ", err)
	}
	//fmt.Printf("\r[+]%v subdomains gathered", countedLines)
	myLogger.Info("%v new subdomain found!", subCount-oldSubCount)

	// Log process completion and elapsed time
	myLogger.Info("assetfinder executed successfully")

	return nil
}

func RunAmass(domain, filePath string, timeout int) error {
	myLogger.Info("Running amass")

	// printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "amass")

	// Show progress
	utils.ShowProgress()

	// Run amass
	output, err := utils.RunCommand("amass", "enum", "-passive", "-timeout", timeout, "-d", domain)
	if err != nil {
		return err
	}

	// Process amass output and apply filter for getting subdomains
	var subdomains []string
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "FQDN") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				subdomain := fields[0]
				subdomains = append(subdomains, subdomain)
			}
		}
	}

	// Join subdomains into a byte slice
	filteredOutput := []byte(strings.Join(subdomains, "\n") + "\n")

	// Write output to specified file
	err = utils.AppendToFile(filePath, filteredOutput)
	if err != nil {
		myLogger.Warning("Error appending to file %s: %v", filePath, err)
	}

	// Count enumareted subdomains
	oldSubCount := subCount
	subCount, err := utils.CountLines(filePath)
	if err != nil {
		myLogger.Warning("Error measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v new subdomain found!", subCount-oldSubCount)

	// Log process completion and elapsed time
	myLogger.Info("amass executed successfully")

	return nil
}

func RunSubkill3r(domain, filePath, wordlist, serverAddr string, workerCount int) error {
	myLogger.Info("Running subkill3r")

	// printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "subkill3r")

	// Show progress
	utils.ShowProgress()

	var filteredResults []string
	results, err := subkill3r.Subkill3r(domain, wordlist, serverAddr, workerCount)
	if err != nil {
		return err
	}

	// Apply filter on gathered results to extract subdomains
	for _, r := range results {
		filteredResults = append(filteredResults, r.Hostname, "\n")
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		myLogger.Error("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	// Iterate through filteredResults and write each []byte to the file
	for _, result := range filteredResults {
		_, err := file.Write([]byte(result))
		if err != nil {
			myLogger.Warning("Error writing to file %s: %v", filePath, err)
		}
	}

	// Count enumareted subdomains
	oldSubCount := subCount
	subCount, err := utils.CountLines(filePath)
	if err != nil {
		log.Printf("Error measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v new subdomain found!", subCount-oldSubCount)

	myLogger.Info("subkill3r executed successfully")

	return nil
}

func InitSubdEnum(cfg PassiveEnum) error {
	modName := "PASSIVE_ENUM"
	myLogger.Info(color.CyanString("%s module initialized\n", modName))

	// FATAL inital foothold for subd enum (can be altered later)
	if err := RunSubfinder(cfg.Domain, cfg.FilePath, cfg.Subfinder.NumOfThreads); err != nil {
		return fmt.Errorf(color.RedString("Error running subfinder for domain %s: %v\n", cfg.Domain, err))
	}

	if err := RunAssetfinder(cfg.Domain, cfg.FilePath); err != nil {
		myLogger.Error("Error running assetfinder for domain %s: %v\n", cfg.Domain, err)
	}

	if err := RunAmass(cfg.Domain, cfg.FilePath, cfg.Amass.Timeout); err != nil {
		myLogger.Error("Error running amass for cfg.Domain %s: %v\n", cfg.Domain, err)
	}

	if cfg.Subkill3r.Wordlist != "none" {
		if err := RunSubkill3r(cfg.Domain, cfg.FilePath, cfg.Subkill3r.Wordlist, cfg.Subkill3r.ServerAddr, cfg.Subkill3r.WorkerCount); err != nil {
			myLogger.Error("Error running subkill3r for cfg.Domain %s: %v", cfg.Domain, err)
			myLogger.Warning("Look for SUBKILL3R_WORDLIST in config file to specify a wordlist\n")
		}
	} else {
		myLogger.Warning("subkill3r is not activated because wordlist is not provided\n")
	}

	// Count total enumerated subdomains
	subCount, err := utils.CountLines(cfg.FilePath)
	if err != nil {
		myLogger.Error("\rError measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v total subdomains gathered", subCount)

	// Removing duplicates: FATAL
	if err := utils.RemoveDuplicatesFromFile(cfg.FilePath); err != nil {
		return fmt.Errorf(color.RedString("%s module failed: error removing duplicates from the file %s: %v ", modName, cfg.FilePath, err))
	}
	myLogger.Info("Removing duplicates from %s", cfg.FilePath)

	// Count unique subdomains
	subCount, err = utils.CountLines(cfg.FilePath)
	if err != nil {
		myLogger.Warning("\rError measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v unique subdomains gathered\n", subCount)

	myLogger.Info(color.CyanString("%s module completed\n", modName))

	return nil
}
