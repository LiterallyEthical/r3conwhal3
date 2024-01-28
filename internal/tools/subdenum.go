package tools

import (
	"fmt"
	"github/literallyethical/pkg/logger"
	"github/literallyethical/pkg/subkill3r"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)


var myLogger logger.Logger
var subCount int


func init() {
	// Init the logger during package initialization
	log, err := logger.NewLogger(0,0,0)
	if err != nil {
		panic(err)
	}
	
	myLogger = log	
}


func RunSubfinder(domain, filePath string) (error) {
	// fmt.Printf("\n[+]Starting subfinder\n")
	myLogger.Info("Running subfinder")
	
	// printing the execution time
	startTime := time.Now()
	defer logElapsedTime(startTime, "subfinder")
	
	// Show progress
	showProgress()

	// Run subfinder
	_ , err := runCommand("subfinder", "-d", domain, "-o", filePath)
	if err != nil {
		return err
	}
	

	// Count enumareted subdomains
	subCount, err = countLines(filePath)
	if err != nil {
		// log.Printf("Error measuring enumerated subdomains: %v ", err) 
		myLogger.Warning("Failed to measure number of gathered subdomains: %v", err)
	}
	// fmt.Printf("\r[+]%v subdomains gathered", countedLines)
	myLogger.Info("%v new subdomain found!", subCount)

	// Log process completion and elapsed time
	// fmt.Printf("\n[+]Subfinder executed successfully")
	myLogger.Info("Subfinder executed successfully")
	return nil	
}

func RunAssetfinder(domain, filePath string) (error) {
	myLogger.Info("Running assetfinder")
	
	// printing the execution time
	startTime := time.Now()
	defer logElapsedTime(startTime, "assetfinder")
	
	// Show progress
	showProgress()

	// Run assetfinder
	output , err := runCommand("assetfinder", "-subs-only", domain)
	if err != nil {
		return err
	}



	
	// Write output to specified file
	err = appendToFile(filePath, output)
	if err != nil {
		//log.Printf("Error appending to file %s: %v", filePath, err)
		myLogger.Warning("Error appending to file %s: %v", filePath, err)
	}

	// Count enumareted subdomains
	
	oldSubCount := subCount
	subCount, err = countLines(filePath)
	if err != nil {
		//log.Printf("Error measuring enumerated subdomains: %v ", err) 
		myLogger.Warning("Error measuring enumerated subdomains: %v ", err)
	}
	//fmt.Printf("\r[+]%v subdomains gathered", countedLines)
	myLogger.Info("%v new subdomain found!", subCount - oldSubCount)

	// Log process completion and elapsed time
	myLogger.Info("assetfinder executed successfully")

	return nil	
}

func RunAmass(domain, filePath string) (error) {
	myLogger.Info("Running amass")

	// printing the execution time
	startTime := time.Now()
	defer logElapsedTime(startTime, "amass")
	
	// Show progress
	showProgress()

	// Run amass
	output , err := runCommand("amass", "enum", "-passive", "-timeout", "1", "-d", domain)
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
	err = appendToFile(filePath, filteredOutput)
	if err != nil {
		myLogger.Warning("Error appending to file %s: %v", filePath, err)
	}

	// Count enumareted subdomains
	oldSubCount := subCount
	subCount, err := countLines(filePath)
	if err != nil {
		myLogger.Warning("Error measuring enumerated subdomains: %v ", err) 
	}
	myLogger.Info("%v new subdomain found!", subCount - oldSubCount)


	// Log process completion and elapsed time
	myLogger.Info("amass executed successfully")

	return nil	
}

func RunSubkill3r(domain, filePath, wordlist, serverAddr string, workerCount int) (error) {
	myLogger.Info("Running subkill3r")
	
	// printing the execution time
	startTime := time.Now()
	defer logElapsedTime(startTime, "subkill3r")

	// Show progress
	showProgress()

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
	subCount, err := countLines(filePath)
	if err != nil {
		log.Printf("Error measuring enumerated subdomains: %v ", err) 
	}
	myLogger.Info("%v new subdomain found!", subCount-oldSubCount)

	myLogger.Info("subkill3r executed successfully")
	
	return nil
}

func RunHTTPX(filePath, dirPath string) error {
	// filter live subdomains
	myLogger.Info("Running httpx")
	liveSubdomains := filepath.Join(dirPath, "live_subdomains.txt")
	
	// printing the execution time
	startTime := time.Now()
	defer logElapsedTime(startTime, "httpx")

	// Show progress
	showProgress()

	_, err := runCommand("httpx", "-fc", "404", "-l", filePath, "-o", liveSubdomains)
	if err != nil {
		return err
	}

	// //Write output to specified file
	// err = appendToFile(liveSubdomains, output)
	// if err != nil {
	// 	return fmt.Errorf("Error appending to file %s: %v", filePath, err)
	// }

	subCount, err := countLines(liveSubdomains)
	if err != nil {
		myLogger.Warning("Failed to measure live subdomains: %v", err)
	}
	myLogger.Info("%v live subdomain found!", subCount)

	// Log process completion and elapsed time
	myLogger.Info("httpx executed successfully")

	return nil
}



func InitSubdEnum(domain, filePath, dirPath, wordlist, serverAddr string, workerCount int) error {
	fmt.Println()
	myLogger.Info(color.CyanString("subdenum module initialized\n"))

	// FATAL inital foothold for subd enum (can be altered later)
	if err := RunSubfinder(domain, filePath); err != nil {
		return fmt.Errorf(color.RedString("Error running subfinder for domain %s: %v\n", domain, err))
	}

	if err := RunAssetfinder(domain, filePath); err != nil {
		myLogger.Error("Error running assetfinder for domain %s: %v\n", domain, err)
	}

	// if err := RunAmass(domain, filePath); err != nil {
	// 	myLogger.Error("Error running amass for domain %s: %v\n", domain, err)	 
	// }

 	if err := RunSubkill3r(domain, filePath, wordlist, serverAddr, workerCount); err != nil {
		myLogger.Error("Error running custom subdomain enumerator for domain %s: %v\n", domain, err)
	 }
 	
	// Count total enumerated subdomains
	subCount, err := countLines(filePath)
	if err != nil {
		myLogger.Error("\rError measuring enumerated subdomains: %v ", err)
	}
	myLogger.Info("%v total subdomains gathered", subCount)
	
	
	// Removing duplicates: FATAL
	if err := RemoveDuplicatesFromFile(filePath); err != nil {
	 	return fmt.Errorf(color.RedString("subenum module failed: error removing duplicates from the file %s: %v ", filePath, err))
	}
	myLogger.Info("Removing duplicates from %s", filePath)
 	
	// Count unique subdomains
	subCount, err = countLines(filePath)
	if err != nil {
		myLogger.Warning("\rError measuring enumerated subdomains: %v ", err) 
	}
	myLogger.Info("%v unique subdomains gathered\n", subCount)

	// Filter live domains: FATAL
	if err := RunHTTPX(filePath, dirPath); err != nil {
		return fmt.Errorf(color.RedString("subenum module failed: error running httpx for %s: %v\n", filePath, err))
	}

	myLogger.Info(color.CyanString("subdenum module completed"))
	
	return nil
}



