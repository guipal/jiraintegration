package main

import (
	"flag"

	"fmt"

	"../jiraUtils"

	"os"
)

func main() {

	var mode, directory, hostUrl string
	var filter int

	flag.StringVar(&mode, "mode", "EXPORT", " [ EXPORT Export tests from Jira | IMPORT Import tests results to Jira | UPDATE Update with HTML report URL to Jenkins]")

	flag.StringVar(&directory, "outputDir", ".", "Target directory for exported tests")

	flag.IntVar(&filter, "filter", 0, "Filter query to retrieve tests")

	flag.StringVar(&hostUrl, "host", "", "Jira server URL")

	flag.Parse()

	if hostUrl == "" {
		fmt.Println("Missing host url, try " + os.Args[0] + " -h")
		return
	}

	jirautils.ExportTests(hostUrl, filter, directory)

}
