package mods

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/fatih/color"
)

type WebOps struct {
	OutDirPath      string
	Gowitness       Gowitness
	FFUF            FFUF
	EnableGowitness bool
	EnableFFUF      bool
}

type Gowitness struct {
	Timeout               int
	ResolutionX           int
	ResolutionY           int
	NumOfThreads          int
	Fullpage              bool
	ScreenshotFilter      bool
	ScreenshotFilterCodes string
}

type FFUF struct {
	NumOfThreads       int
	Maxtime            int
	Rate               int
	Timeout            int
	Wordlist           string
	MatchHTTPCode      string
	FilterResponseSize string
	OutputFormat       string
	Output             string
	SF                 bool
	SE                 bool
}

func RunGowitness(outdirPath string, timeout, resolutionX, resolutionY, numOfThreads int, fullpage, screenshotFilter bool, screenshotFilterCodes string) error {

	myLogger.Info("Running gowitness")

	// Printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "gowitness")

	// Show progress
	utils.ShowProgress()

	// Run RunGowitness
	screenshotPath := filepath.Join(outdirPath, "screenshots/")
	filename := filepath.Join(outdirPath, "live_subdomains.txt")

	// Prepare the command arguments
	cmdArgs := []string{
		"file",
		"-f", filename,
		"--screenshot-path", screenshotPath,
		"--disable-db",
		"--timeout", fmt.Sprint(timeout),
		"--resolution-x", fmt.Sprint(resolutionX),
		"--resolution-y", fmt.Sprint(resolutionY),
		"--threads", fmt.Sprint(numOfThreads),
	}

	if fullpage {
		cmdArgs = append(cmdArgs, "--fullpage")
	}

	if screenshotFilter {
		cmdArgs = append(cmdArgs, "--screenshot-filter", screenshotFilterCodes)
	}

	// Convert cmdArgs to []interface{} for RunCommand
	interfaceArgs := make([]interface{}, len(cmdArgs))
	for i, arg := range cmdArgs {
		interfaceArgs[i] = arg
	}

	// Run gowitness
	_, err := utils.RunCommand("gowitness", interfaceArgs...)
	if err != nil {
		return err
	}

	return nil
}

func RunFFUF(numOfThreads, maxtime, rate, timeout int, outDirPath, wordlist, matchHTTPCode, filterResponseSize, outputFormat, output string, SF, SE bool) error {

	myLogger.Info("Running ffuf")

	// Printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "ffuf")

	// Show progress
	utils.ShowProgress()

	// wordlist
	src := filepath.Join(outDirPath, "live_subdomains.txt")
	domainWordlist := filepath.Join(".tmp", "live_subdomains.txt")

	if err := utils.CopyFile(src, domainWordlist); err != nil {
		return fmt.Errorf("Failed to copy file: %s", err)
	}

	// output dst
	dstDir := filepath.Join(outDirPath, "web_ops")
	if err := os.Mkdir(dstDir, 0755); err != nil {
		return err
	}
	outFile := filepath.Join(dstDir, fmt.Sprintf("%s.%s", output, outputFormat))

	// Prepare the command arguments
	cmdArgs := []string{
		"-u", "FUZZDOMAIN/FUZZDIR",
		"-w", fmt.Sprintf("%s:FUZZDOMAIN", domainWordlist),
		"-w", fmt.Sprintf("%s:FUZZDIR", wordlist),
		"-mc", matchHTTPCode,
		"-fs", filterResponseSize,
		"-maxtime", fmt.Sprint(maxtime),
		"-of", outputFormat,
		"-o", outFile,
		"-timeout", fmt.Sprint(timeout),
		"-t", fmt.Sprint(numOfThreads),
		"-rate", fmt.Sprint(rate),
		"-s",
	}

	if SF {
		cmdArgs = append(cmdArgs, "-sf")
	}

	if SE {
		cmdArgs = append(cmdArgs, "-se")
	}

	// Convert cmdArgs to []interface{} for RunCommand
	interfaceArgs := make([]interface{}, len(cmdArgs))
	for i, arg := range cmdArgs {
		interfaceArgs[i] = arg
	}

	// Run FFUF
	_, err := utils.RunCommand("ffuf", interfaceArgs...)
	if err != nil {
		return err
	}

	return nil
}

func InitWebOps(cfg WebOps) error {
	modName := "WEB_OPS"
	myLogger.Info(color.MagentaString("%s module initialized\n", modName))

	if cfg.EnableGowitness {

		// Web screenshoting
		myLogger.Info(color.MagentaString("WEB_SCREENSHOTING is activated"))
		if err := RunGowitness(cfg.OutDirPath, cfg.Gowitness.Timeout, cfg.Gowitness.ResolutionX, cfg.Gowitness.ResolutionY, cfg.Gowitness.NumOfThreads, cfg.Gowitness.Fullpage, cfg.Gowitness.ScreenshotFilter, cfg.Gowitness.ScreenshotFilterCodes); err != nil {
			return fmt.Errorf(color.RedString("Error running gowitness: %v", err))
		}
	}

	if cfg.EnableFFUF {
		// Directory fuzzing
		myLogger.Info(color.MagentaString("DIRECTORY_FUZZING is activated"))
		if err := RunFFUF(cfg.FFUF.NumOfThreads, cfg.FFUF.Maxtime, cfg.FFUF.Rate, cfg.FFUF.Timeout, cfg.OutDirPath, cfg.FFUF.Wordlist, cfg.FFUF.MatchHTTPCode, cfg.FFUF.FilterResponseSize, cfg.FFUF.OutputFormat, cfg.FFUF.Output, cfg.FFUF.SF, cfg.FFUF.SE); err != nil {
			return fmt.Errorf(color.RedString("Error running FFUF: %v", err))
		}
	}

	myLogger.Info(color.MagentaString("%s module completed\n", modName))

	return nil
}
