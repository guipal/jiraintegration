package main

import (
	"flag"

	"fmt"

	"os"

	"github.com/guipal/jiraintegration/jirautils"

	"github.com/howeyc/gopass"
)

func main() {

	var directory, hostUrl, user, password, keys, resultPath string
	var exportResult, importTest, executeTest, executeRemote bool
	var filter int

	flag.BoolVar(&exportResult, "export", false, "Export cucumber resutlts to Jira")
	flag.BoolVar(&importTest, "download_test", false, "Download test Scenarios from Jira")
	flag.BoolVar(&executeTest, "execute", true, "Execute test Scenarios from $featuresDir")
	flag.BoolVar(&executeRemote, "executeRemote", false, "Download, Execute & Upload test from/to $host server")

	flag.StringVar(&directory, "featuresDir", "./features", "Target directory for downloaded tests")

	flag.IntVar(&filter, "filter", 0, "Filter query to retrieve tests")

	flag.StringVar(&hostUrl, "host", "", "Jira server URL")
	flag.StringVar(&user, "user", "", "Jira server user")
	flag.StringVar(&password, "password", "", "Jira server password")
	flag.StringVar(&keys, "keys", "", " list of Jira-keys of the tests separated by ';'")
	flag.StringVar(&resultPath, "resultFile", "./results.json", " Path to cucumber result file")

	flag.Parse()

	if (importTest || exportResult || executeRemote) && hostUrl == "" {
		fmt.Println("Missing host url, try " + os.Args[0] + " -h")
		return
	}

	if (importTest || exportResult || executeRemote) && password == "" {
		fmt.Printf("Password: ")
		pass, _ := gopass.GetPasswd()
		password = string(pass)
	}

	if executeRemote {
		err := jirautils.ExecuteTestSet(hostUrl, filter, directory, user, password, keys, resultPath)
		if err != nil {
			fmt.Println(err)
		}
	} else {

		if importTest {
			err := jirautils.DownloadTests(hostUrl, filter, directory, user, password, keys)
			if err != nil {
				fmt.Println(err)
			}
		}
		if executeTest {
			err := jirautils.ExecuteCucumberTest(resultPath, directory)
			if err != nil {
				fmt.Println(err)
			}
		}

		if exportResult {
			err := jirautils.ExportTestExecution(hostUrl, resultPath, user, password)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

}
