package main

import (
	"flag"

	"fmt"

	"os"

	"github.com/guipal/jiraintegration/jirautils"

	"github.com/howeyc/gopass"
)

func main() {

	var mode, directory, hostUrl, user, password, keys, resultPath string
	var filter int

	flag.StringVar(&mode, "mode", "EXPORT", " [ EXPORT Export tests from Jira | IMPORT Import tests results to Jira | UPDATE Update with HTML report URL to Jenkins | EXECUTE Execute tests and upload results]")

	flag.StringVar(&directory, "outputDir", "./features", "Target directory for exported tests")

	flag.IntVar(&filter, "filter", 0, "Filter query to retrieve tests")

	flag.StringVar(&hostUrl, "host", "", "Jira server URL")
	flag.StringVar(&user, "user", "", "Jira server user")
	flag.StringVar(&password, "password", "", "Jira server password")
	flag.StringVar(&keys, "keys", "", " list of Jira-keys of the tests separated by ';'")
	flag.StringVar(&resultPath, "resultFile", "./results.json", " Path to cucumber result file")

	flag.Parse()

	if hostUrl == "" {
		fmt.Println("Missing host url, try " + os.Args[0] + " -h")
		return
	}

	if password == "" {
		fmt.Printf("Password: ")
		pass, _ := gopass.GetPasswd()
		password = string(pass)
	}

	switch mode {

	case "EXPORT":
		jirautils.ExportTests(hostUrl, filter, directory, user, password, keys)
	case "IMPORT":
		jirautils.ImportTestsExecution(hostUrl, resultPath, user, password)
	case "EXECUTE":
		err := jirautils.ExecuteTestSet(hostUrl, filter, "./features", user, password, keys, resultPath)
		if err != nil {
			fmt.Println(err)
		}

	}

}
