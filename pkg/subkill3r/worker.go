package subkill3r

type empty struct{}

func Worker(tracker chan empty, fqdns chan string, gather chan []Result, serverAddr string) {
	for fqdn := range fqdns {
		results := Lookup(fqdn, serverAddr)
		if len(results) > 0  {
			gather <- results
		}
	}
	var e empty
	tracker <- e // send signals to tracker to inform the job has done
}