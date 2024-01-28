package subkill3r

import (
	"bufio"
	"os"
)

// Result represents the result of a subdomain lookup.
type Result struct {
	IPAdress string
	Hostname string
}

// Subkill3r performs subdomain enumeration and returns the results.
func Subkill3r(domain, wordlist, serverAddr string, workerCount int) ([]Result, error) {
	var results []Result
	fqdns := make(chan string, workerCount)
	gather := make(chan []Result)
	tracker := make(chan empty)

	fh, err := os.Open(wordlist)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)

	// Initializing the Worker goroutines
	for i := 0; i < workerCount; i++ {
		go Worker(tracker, fqdns, gather, serverAddr)
	}

	// Populating subdomains via reading from file
	for scanner.Scan() {
		fqdns <- formatFQDN(scanner.Text(), domain)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Gathering and returning the results
	go func() {
		for r := range gather {
			results = append(results, r...)
		}
		var e empty
		tracker <- e
	}()

	close(fqdns) // No longer data will be sent to channel
	for i := 0; i < workerCount; i++ {
		<-tracker
	}
	close(gather)
	<-tracker

	return results, nil
}

// formatFQDN formats the fully qualified domain name.
func formatFQDN(subdomain, domain string) string {
	return subdomain + "." + domain
}