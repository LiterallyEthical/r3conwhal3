package utils

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/LiterallyEthical/r3conwhal3/pkg/logger"
	"github.com/fatih/color"
)

var myLogger logger.Logger

func init() {
	// Init the logger during package initialization
	log, err := logger.NewLogger(0, 0, 0)
	if err != nil {
		panic(err)
	}

	myLogger = log
}

func RunCommand(command string, args ...interface{}) ([]byte, error) {

	var strArgs []string
	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			strArgs = append(strArgs, strconv.Itoa(v))
		case string:
			strArgs = append(strArgs, v)
		case bool:
			strArgs = append(strArgs, strconv.FormatBool(v))
		default:
			return nil, fmt.Errorf("unsupported argument type: %T", v)
		}
	}

	cmd := exec.Command(command, strArgs...)

	// Capture command output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("\n[-]error running %s: %v\n%s", command, err, output)
	}

	return output, nil
}

func AppendToFile(filePath string, content []byte) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening or creating file %s: %v", filePath, err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error appending to file %s: %v", filePath, err)
	}

	return nil
}

func ShowProgress() {
	// Simulate progress with a simple animation
	animation := []string{".", "..", "...", "...."}

	//Clear the animation line
	defer fmt.Print("\r", "")

	for _, frame := range animation {
		fmt.Print("\r", frame)
		time.Sleep(500 * time.Millisecond)
	}
}

func LogElapsedTime(startTime time.Time, operation string) {
	elapsedTime := time.Since(startTime)
	myLogger.Info("%s completed in %s\n", operation, elapsedTime)
}

func CheckInstallations(tools []string) error {
	// fmt.Printf("[+]Start checking required tools\n")
	myLogger.Info(color.CyanString("Checking required tools\n"))

	ShowProgress()

	//versionRegex := regexp.MustCompile(`(\d+\.\d+\.\d+)`)

	for _, tool := range tools {
		cmd := exec.Command("which", tool)
		output, err := cmd.CombinedOutput()

		if err != nil {
			// Check if the error is an "ExitError" and if the exit code is 1
			exitErr, ok := err.(*exec.ExitError)
			if ok && exitErr.ExitCode() == 1 {
				return fmt.Errorf("\n%s is not installed or not in the system's PATH", tool)
			}
			// Return the general error if it's not an "ExitError" or if the exit code is not 127
			return fmt.Errorf("\nerror running %s: %v", tool, err)
		}

		// Check if the output is empty, that indicates tool is not installed
		if strings.TrimSpace(string(output)) == "" {
			return fmt.Errorf("\n%s is not in the system's PATH", tool)
		}

		// fmt.Printf("\n[+]%s is installed", tool)
		myLogger.Info("%s is installed", tool)
	}

	return nil
}

func CountLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

// Remove duplicates from the given file
func RemoveDuplicatesFromFile(filename string) error {
	// Open the file with read and write permissions
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a map to store unique lines
	uniqueLines := make(map[string]bool)

	// Scan the file line by line
	for scanner.Scan() {
		// Read the current line
		line := scanner.Text()
		uniqueLines[line] = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Move the file cursor to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	// Truncate the file to remove existing content
	if err := file.Truncate(0); err != nil {
		return err
	}

	// Writes to file without duplicates
	writer := bufio.NewWriter(file)

	for line := range uniqueLines {
		// Write each unique line to the file
		fmt.Fprintln(writer, line)
	}

	// ensure all data is written to the file
	writer.Flush()

	return nil
}

func CreateDir(dirName, domain string) (string, error) {
	// Check if the directory already exists
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			return "", fmt.Errorf("error creating base directory: %v", err)
		}
	}

	// Extracting the target domain name
	domainPrefix := strings.Split(domain, ".")[0]

	// Get the current timestamp
	timestamp := time.Now().Format("2006-01-02-15:04")

	// Combine domain prefix with timestamp
	subdirName := fmt.Sprintf("%s_%s", domainPrefix, timestamp)

	// Create the subdirectory
	subdirPath := filepath.Join(dirName, subdirName)
	if err := os.MkdirAll(subdirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating subdirectory: %v", err)
	}

	return subdirPath, nil
}

// print banner in ascii art format
func Banner(bannerPath string) {
	b, err := os.ReadFile(bannerPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(color.CyanString(string(b)))
}

// ExtractEmbeddedFileToTempDir reads an embedded file and writes it to a temporary directory named ".tmp", returning the path to the newly created temporary file.
func ExtractEmbeddedFileToTempDir(docFS embed.FS, embeddedFilePath, tempFileName string) (string, error) {
	tempDir := ".tmp" // Hardcoded temporary directory

	// Ensure the temporary directory exists
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %v", err)
	}

	// Read the embedded file
	data, err := docFS.ReadFile(embeddedFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded file '%s': %v", embeddedFilePath, err)
	}

	// Write to a temporary file within the specified temporary directory
	tmpFilePath := filepath.Join(tempDir, tempFileName)
	if err := os.WriteFile(tmpFilePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write to temporary file '%s': %v", tmpFilePath, err)
	}

	return tmpFilePath, nil
}

// It search for existence of specificFiles in the given directory and merge them to a new file
func MergeFiles(pathToDir, outFileName string, specificFiles []string) error {
	var mergedContent []byte

	// Iterate over the list of specific values
	for _, fileName := range specificFiles {
		filePath := filepath.Join(pathToDir, fileName)

		// Check if the file exists
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				// The file does not exist, so it needs to be handled
				myLogger.Warning("File does not exist: %s\n", filePath)
				continue // Skip this file
			} else {
				return err
			}
		}

		// Read the content of the file
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		// Append the content of current file to the mergedContent
		mergedContent = append(mergedContent, content...)
	}

	// Construct the path for outFileName
	outPath := filepath.Join(pathToDir, outFileName)

	// Write the merged content to the specified file
	err := os.WriteFile(outPath, mergedContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
