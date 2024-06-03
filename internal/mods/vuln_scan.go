package mods

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/LiterallyEthical/r3conwhal3/internal/utils"
	"github.com/fatih/color"
)

type VulnScan struct {
	OutdirPath  string
	Subzy       Subzy
	EnableSubzy bool
}

type Subzy struct {
	Concurrency int
	Timeout     int
	HideFails   bool
	HTTPS       bool
	VerifySSL   bool
	Vuln        bool
}

type SubzyOutput struct {
	Subdomain     string `json:"subdomain,omitempty"`
	Status        string `json:"status,omitempty"`
	Engine        string `json:"engine,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	Discussion    string `json:"discussion,omitempty"`
}

func RunSubzy(outdirPath string, concurrency, timeout int, hideFails, HTTPS, verifySSL, vuln bool) error {

	myLogger.Info("Running subzy")

	// Printing the execution time
	startTime := time.Now()
	defer utils.LogElapsedTime(startTime, "subzy")

	// Show progress
	utils.ShowProgress()

	// Prepare the command arguments
	outFolder := filepath.Join(outdirPath, "vuln_scan")
	if err := os.Mkdir(outFolder, 0755); err != nil {
		return fmt.Errorf("Error while creating the directory vuln_scan: %v", err)
	}

	outFile := filepath.Join(outFolder, "subdomain_takeover_scan.json")
	targets := filepath.Join(outdirPath, "live_subdomains.txt")

	cmdArgs := []string{
		"run",
		"--targets", targets,
		"--output", outFile,
		"--concurrency", fmt.Sprint(concurrency),
		"--timeout", fmt.Sprint(timeout),
	}

	if hideFails {
		cmdArgs = append(cmdArgs, "--hide_fails")
	}

	if HTTPS {
		cmdArgs = append(cmdArgs, "--https")
	}

	if verifySSL {
		cmdArgs = append(cmdArgs, "--verify_ssl")
	}

	if vuln {
		cmdArgs = append(cmdArgs, "--vuln")
	}

	// Convert cmdArgs to []interface{} for RunCommand
	interfaceArgs := make([]interface{}, len(cmdArgs))
	for i, arg := range cmdArgs {
		interfaceArgs[i] = arg
	}

	// Run gowitness
	_, err := utils.RunCommand("subzy", interfaceArgs...)
	if err != nil {
		return fmt.Errorf("Error while running command subzy: %v", err)
	}

	// Read the JSON file
	file, err := os.Open(outFile)
	if err != nil {
		return fmt.Errorf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Read the file's content
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Failed to read file: %s", err)
	}

	// Parse the JSON data
	var out []SubzyOutput
	err = json.Unmarshal(byteValue, &out)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal JSON: %s", err)
	}

	// Log the subzy output
	anyVuln := false
	subzyLabels := [3]string{"VULNERABLE", "DISCUSSION", "DOCUMENTATION"}
	for _, o := range out {
		if o.Status == "vulnerable" {
			myLogger.Info("The following subdomain flagged as vulnerable to subdomain takeover")
			myLogger.Info("[ %s ]  -  %s [ %s ]", color.RedString(subzyLabels[0]), o.Subdomain, color.RedString(o.Engine))
			myLogger.Info("[ %s ]  -  %s", color.YellowString(subzyLabels[1]), o.Discussion)
			myLogger.Info("[ %s ]  -  %s", color.CyanString(subzyLabels[2]), o.Documentation)
			anyVuln = true
		}
	}

	if !anyVuln {
		myLogger.Info("Target subdomains are not vulnerable to subdomain takeover")
	}

	myLogger.Info("Subzy executed successfully")
	return nil
}

func InitVulnScan(cfg VulnScan) error {
	modName := "VULN_SCAN"
	myLogger.Info(color.YellowString("%s module initialized\n", modName))

	if cfg.EnableSubzy {
		if err := RunSubzy(cfg.OutdirPath, cfg.Subzy.Concurrency, cfg.Subzy.Timeout, cfg.Subzy.HideFails, cfg.Subzy.HTTPS, cfg.Subzy.VerifySSL, cfg.Subzy.Vuln); err != nil {
			return fmt.Errorf(color.RedString("Error running subzy: %v\n", err))
		}
	}

	myLogger.Info(color.YellowString("%s module completed\n", modName))

	return nil
}
