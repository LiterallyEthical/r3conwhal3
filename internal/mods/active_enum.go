package mods

import (
	"fmt"
	"time"

	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/fatih/color"
)

func RunPureDNS(domain, filePath, wordlist, resolvers string, numOfThreads int) error {

	myLogger.Info("Running puredns")

	// printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "puredns")

	// Show progress
	utils.ShowProgress()

	// Run puredns
	_, err := utils.RunCommand("puredns", "bruteforce", wordlist, domain, "--resolvers", resolvers, "--write", filePath, "--threads", numOfThreads)
	if err != nil {
		return err
	}

	// Count resolved subdomains
	subCount, err = utils.CountLines(filePath)
	if err != nil {
		myLogger.Warning("Failed to measure number of resolved subdomains: %v", err)
	}

	myLogger.Info("%v new subdomain found!", subCount)
	myLogger.Info("Puredns executed successfuully")

	return nil
}

func InitActiveSubdEnum(outDirPath, domain, filePath, wordlist, resolvers string, numOfThreads int) error {
	modName := "ACTIVE_ENUM"
	myLogger.Info(color.CyanString("%s module initialized\n", modName))

	// FATAL inital foothold for this module(can be altered later)
	if err := RunPureDNS(domain, filePath, wordlist, resolvers, numOfThreads); err != nil {
		return fmt.Errorf(color.RedString("Error running puredns for domain %s: %v\n", domain, err))
	}

	// Count total enumerated subdomains
	subCount, err := utils.CountLines(filePath)
	if err != nil {
		myLogger.Error("\rError measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v total subdomains gathered", subCount)

	// Removing duplicates: FATAL
	if err := utils.RemoveDuplicatesFromFile(filePath); err != nil {
		return fmt.Errorf(color.RedString("%s module failed: error removing duplicates from the file %s: %v ", modName, filePath, err))
	}
	myLogger.Info("Removing duplicates from %s", filePath)

	// Count unique subdomains
	subCount, err = utils.CountLines(filePath)
	if err != nil {
		myLogger.Warning("\rError measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v unique subdomains gathered\n", subCount)
	myLogger.Info(color.CyanString("%s module completed", modName))

	return nil
}
