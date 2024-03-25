package mods

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/fatih/color"
)

type ActiveEnum struct {
	PureDNS        PureDNS
	Gotator        Gotator
	OutDirPath     string
	SpecifiedFiles []string
}

type PureDNS struct {
	Domain       string
	Wordlist     string
	Resolvers    string
	NumOfThreads int
}

type Gotator struct {
	Sublist      string
	Permlist     string
	Depth        int
	Numbers      int
	NumOfThreads int
	Mindup       bool
	Adv          bool
	Md           bool
}

func RunPureDNS(mode, domain, outDirPath, wordlist, resolvers string, numOfThreads int) error {

	myLogger.Info("Running puredns")

	// printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "puredns")

	// Show progress
	utils.ShowProgress()

	// Run puredns
	switch mode {
	case "bruteforce":
		outDir := filepath.Join(outDirPath, "active_enum_subdomains.txt")
		_, err := utils.RunCommand("puredns", "bruteforce", wordlist, domain, "--resolvers", resolvers, "--write", outDir, "--threads", numOfThreads)
		if err != nil {
			return err
		}

		// Count resolved subdomains
		subCount, err := utils.CountLines(outDir)
		if err != nil {
			myLogger.Warning("Failed to measure number of resolved subdomains: %v", err)
		}

		myLogger.Info("%v new subdomain found!", subCount)

	case "resolve":
		// Ensure the temporary directory exists
		tempDir := ".tmp"
		if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create temporary directory: %v", err)
		}
		outDir := filepath.Join(outDirPath, "resolved_subs.txt")
		permutatedSubs := filepath.Join(tempDir, "permutated_subs.txt")

		_, err := utils.RunCommand("puredns", "resolve", permutatedSubs, "--resolvers", resolvers, "--write", outDir, "--threads", numOfThreads)
		if err != nil {
			return err
		}

		// Count resolved subdomains
		subCount, err := utils.CountLines(outDir)
		if err != nil {
			myLogger.Warning("Failed to measure number of resolved subdomains: %v", err)
		}

		myLogger.Info("%v new subdomain found!", subCount)

	default:
		return fmt.Errorf("Unknown mode for puredns: %v", mode)
	}

	myLogger.Info("Puredns executed successfully\n")

	return nil
}

func RunGotator(sublist, permlist string, depth, numbers, numOfThreads int, mindup, adv, md bool) error {

	myLogger.Info("Running gotator")

	// printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "gotator")

	// Show progress
	utils.ShowProgress()

	// Prepare the command arguments
	cmdArgs := []string{
		"-sub", sublist,
		"-perm", permlist,
		"-depth", fmt.Sprint(depth),
		"-numbers", fmt.Sprint(numbers),
		"-t", fmt.Sprint(numOfThreads),
		"-silent",
	}

	// Add flags based on their boolean value
	if mindup {
		cmdArgs = append(cmdArgs, "-mindup")
	}
	if adv {
		cmdArgs = append(cmdArgs, "-adv")
	}
	if md {
		cmdArgs = append(cmdArgs, "-md")
	}

	// Convert cmdArgs to []interface{} for RunCommand
	interfaceArgs := make([]interface{}, len(cmdArgs))
	for i, arg := range cmdArgs {
		interfaceArgs[i] = arg
	}

	// Run gotator
	output, err := utils.RunCommand("gotator", interfaceArgs...)
	if err != nil {
		return err
	}

	// Write output to specified file
	filePath := filepath.Join(".tmp", "permutated_subs.txt")
	err = utils.AppendToFile(filePath, output)
	if err != nil {
		return fmt.Errorf("Error appending to file %s: %v", filePath, err)
	}

	lineCount, err := utils.CountLines(filePath)
	if err != nil {
		myLogger.Warning("Error measuring permutated subdomains: %v ", err)
	}

	myLogger.Info("%v permutations generated!", lineCount)
	myLogger.Info("Gotator executed successfully\n")

	return nil
}

func RunMergeFiles(outDirPath, outFileName string, specifiedFiles []string) error {
	myLogger.Info("Starting to merge all subdomains previously gathered")
	outFilePath := filepath.Join(outDirPath, outFileName)

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
	myLogger.Info("%v subdomains gathered", subCount)

	// Removing duplicates: FATAL
	if err := utils.RemoveDuplicatesFromFile(outFilePath); err != nil {
		return fmt.Errorf(color.RedString("ACTIVE_ENUM module failed: error removing duplicates from the file %s: %v ", outFilePath, err))
	}
	myLogger.Info("Removing duplicates from %s", outFilePath)

	// Count unique subdomains
	subCount, err = utils.CountLines(outFilePath)
	if err != nil {
		myLogger.Warning("\rError measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v unique subdomains gathered\n", subCount)

	return nil
}

func InitActiveSubdEnum(cfg ActiveEnum) error {
	modName := "ACTIVE_ENUM"
	myLogger.Info(color.CyanString("%s module initialized\n", modName))

	// FATAL inital foothold for this module(can be altered later)
	myLogger.Info(color.RedString("DNS_BRUTEFORCE is activated"))
	if err := RunPureDNS("bruteforce", cfg.PureDNS.Domain, cfg.OutDirPath, cfg.PureDNS.Wordlist, cfg.PureDNS.Resolvers, cfg.PureDNS.NumOfThreads); err != nil {
		return fmt.Errorf(color.RedString("Error running puredns for domain %s: %v\n", cfg.PureDNS.Domain, err))
	}

	// Merge all subdomain files previously gathered
	myLogger.Info(color.RedString("MERGE_FILES is activated"))
	if err := RunMergeFiles(cfg.OutDirPath, "all_subdomains.txt", cfg.SpecifiedFiles); err != nil {
		return fmt.Errorf(color.RedString("Error running merge files to %v", cfg.OutDirPath))
	}

	// DNS permutation
	myLogger.Info(color.RedString("DNS_PERMUTATION is activated"))
	if err := RunGotator(cfg.Gotator.Sublist, cfg.Gotator.Permlist, cfg.Gotator.Depth, cfg.Gotator.Numbers, cfg.Gotator.NumOfThreads, cfg.Gotator.Mindup, cfg.Gotator.Adv, cfg.Gotator.Md); err != nil {
		return fmt.Errorf(color.RedString("Error running gotator for domain %s: %v\n", cfg.PureDNS.Domain, err))
	}

	// DNS Resolving
	myLogger.Info(color.RedString("DNS_RESOLVE is activated"))
	if err := RunPureDNS("resolve", cfg.PureDNS.Domain, cfg.OutDirPath, cfg.PureDNS.Wordlist, cfg.PureDNS.Resolvers, cfg.PureDNS.NumOfThreads); err != nil {
		return fmt.Errorf(color.RedString("Error running puredns for domain %s: %v\n", cfg.PureDNS.Domain, err))
	}

	myLogger.Info(color.CyanString("%s module completed\n", modName))

	return nil
}
