package main

import (
	"flag"

	"fmt"

	"os"

	"github.com/guipal/jiraintegration/jirautils"

	"github.com/howeyc/gopass"

	"io/ioutil"

	"strings"
)

func main() {

	var directory, hostUrl, user, password, passwordFile, keys, resultPath, format string
	var exportResult, importTest, executeTest, executeRemote bool
	var filter int

	flag.BoolVar(&exportResult, "export", false, "Export cucumber results to Jira")
	flag.BoolVar(&importTest, "download_test", false, "Download test Scenarios from Jira")
	flag.BoolVar(&executeTest, "execute", false, "Execute test Scenarios from $featuresDir")
	flag.BoolVar(&executeRemote, "executeRemote", false, "Download, Execute & Upload test from/to $host server")

	flag.StringVar(&directory, "featuresDir", "./features", "Target directory for downloaded tests")

	flag.IntVar(&filter, "filter", 0, "Filter query to retrieve tests")

	flag.StringVar(&hostUrl, "host", "", "Jira server URL")
	flag.StringVar(&user, "user", "", "Jira server user")
	flag.StringVar(&format, "format", "", "Cucumber test result format")
	flag.StringVar(&password, "password", "", "Jira server password")
	flag.StringVar(&passwordFile, "passwordFile", "", "Jira server password file")
	flag.StringVar(&keys, "keys", "", " list of Jira-keys of the tests separated by ';'")
	flag.StringVar(&resultPath, "resultFile", "", " Path to cucumber result file")

	flag.Parse()

	args := flag.Args()

	if (importTest || exportResult || executeRemote) && hostUrl == "" {
		fmt.Println("Missing host url, try " + os.Args[0] + " -h")
		return
	}

	if !importTest && !exportResult && !executeRemote && !executeTest {
		executeTest = true
	}

	if (importTest || exportResult || executeRemote) && password == "" && passwordFile == "" {
		fmt.Printf("Password: ")
		pass, _ := gopass.GetPasswd()
		password = string(pass)
	} else {
		if passwordFile != "" {
			pass, _ := ioutil.ReadFile(passwordFile)
			password = strings.TrimSpace(string(pass))
		}

	}

	if executeRemote {
		if resultPath == "" {
			resultPath = "./results.json"
		}

		err := jirautils.ExecuteTestSet(hostUrl, filter, directory, user, password, keys, resultPath)
		if err != nil {
			fmt.Println(err)
		}
	} else {

		if importTest {
			err := jirautils.DownloadTests(hostUrl, filter, directory, user, password, keys)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if executeTest {
			if directory == "./features" && len(args) > 0 {
				if exists, _ := jirautils.Exists(string(args[0])); exists {
					directory = args[0]
				}
			}

			err := jirautils.ExecuteCucumberTest(format, resultPath, directory)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if exportResult {
			if resultPath == "" {
				fmt.Println("No file provided")
				os.Exit(1)
			}
			err := jirautils.ExportTestExecution(hostUrl, resultPath, user, password)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	}

}
