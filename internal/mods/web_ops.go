package mods

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/fatih/color"
)

type WebOps struct {
	OutDirPath string
	Gowitness  Gowitness
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

	myLogger.Debug("gowitness CMD: ", cmdArgs)

	// Run gotator
	_, err := utils.RunCommand("gowitness", interfaceArgs...)
	if err != nil {
		return err
	}

	return nil
}

func InitWebOps(cfg WebOps) error {
	modName := "WEB_OPS"
	myLogger.Info(color.MagentaString("%s module initialized\n", modName))

	// Web screenshoting
	myLogger.Info(color.MagentaString("WEB_SCREENSHOTING is activated"))
	if err := RunGowitness(cfg.OutDirPath, cfg.Gowitness.Timeout, cfg.Gowitness.ResolutionX, cfg.Gowitness.ResolutionY, cfg.Gowitness.NumOfThreads, cfg.Gowitness.Fullpage, cfg.Gowitness.ScreenshotFilter, cfg.Gowitness.ScreenshotFilterCodes); err != nil {
		return fmt.Errorf(color.RedString("Error running gowitness: %s", err))
	}

	myLogger.Info(color.MagentaString("%s module completed\n", modName))

	return nil
}
