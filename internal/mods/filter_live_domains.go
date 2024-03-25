package mods

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/fatih/color"
)

func RunHTTPX(filePath, outDirPath string) error {
	// filter live subdomains
	myLogger.Info("Running httpx")
	liveSubdomains := filepath.Join(outDirPath, "live_subdomains.txt")

	// printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "httpx")

	// Show progress
	utils.ShowProgress()

	_, err := utils.RunCommand("httpx", "-l", filePath, "-o", liveSubdomains)
	if err != nil {
		return err
	}

	subCount, err := utils.CountLines(liveSubdomains)
	if err != nil {
		myLogger.Warning("Failed to measure live subdomains: %v", err)
	}
	myLogger.Info("%v live subdomain found!", subCount)

	// Log process completion and elapsed time
	myLogger.Info("httpx executed successfully")

	return nil
}

func InitFilterLiveDomains(outDirPath string) error {

	// Check if the ACTIVE_SUBD_ENUM is run
	if _, err := os.Stat("all_subdomains.txt"); err != nil {
		if os.IsNotExist(err) {
			outFileName := "all_subdomains.txt"
			specifiedFiles := []string{"active_enum_subdomains.txt", "passive_enum_subdomains.txt"}

			// Merge all subdomain files previously gathered
			if err := RunMergeFiles(outDirPath, outFileName, specifiedFiles); err != nil {
				return fmt.Errorf(color.RedString("Error running merge files to %v", outDirPath))
			}
		}
	}

	modName := "FILTER_LIVE_DOMAINS"
	outFileName := "ultimate_subdomains.txt"
	outFilePath := filepath.Join(outDirPath, outFileName)
	specifiedFiles := []string{"all_subdomains.txt", "resolved_subs.txt"}

	myLogger.Info(color.CyanString("%s module initialized", modName))

	// Merge all gathered subdomain files to a single file
	if err := utils.MergeFiles(outDirPath, outFileName, specifiedFiles); err != nil {
		return fmt.Errorf(color.RedString("Error while merging %v to :%v", specifiedFiles, outFileName))
	}
	myLogger.Info("Merge all subdomains operation is successfull!")

	// Count total enumerated subdomains
	subCount, err := utils.CountLines(outFilePath)
	if err != nil {
		myLogger.Error("\rError measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v total subdomains gathered", subCount)

	// Removing duplicates: FATAL
	if err := utils.RemoveDuplicatesFromFile(outFilePath); err != nil {
		return fmt.Errorf(color.RedString("%s module failed: error removing duplicates from the file %s: %v ", modName, outFilePath, err))
	}
	myLogger.Info("Removing duplicates from %s", outFilePath)

	// Count unique subdomains
	subCount, err = utils.CountLines(outFilePath)
	if err != nil {
		myLogger.Warning("\rError measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v total unique subdomains gathered\n", subCount)

	// Filter live subdomains
	if err := RunHTTPX(outFilePath, outDirPath); err != nil {
		return fmt.Errorf(color.RedString("%s module failed: error running httpx for %s: %v\n", modName, outFilePath, err))
	}

	myLogger.Info(color.CyanString("FILTER_LIVE_DOMAINS module completed"))

	return nil
}
