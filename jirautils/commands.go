// Package jiraUtils contains utility functions for working with strings.
package jirautils

import (
	"os"

	"os/exec"

	"net/http"

	"strconv"

	"io/ioutil"

	"errors"

	"bytes"

	"time"

	"strings"
)

// DownloadTests Download selected tests to provided output directory
func DownloadTests(host string, filter int, outputDirectory, user, password, keys string) error {
	os.MkdirAll(outputDirectory, os.ModePerm)
	os.Chdir(outputDirectory)
	defer os.Chdir("../")
	reqUrl := host + "/rest/raven/1.0/export/test?fz=true"
	if filter != 0 {
		reqUrl = reqUrl + "&filter=" + strconv.Itoa(filter)
	}
	if keys != "" {
		reqUrl = reqUrl + "&keys=" + keys
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl, nil)
	req.SetBasicAuth(user, password)
	// ...
	resp, err := client.Do(req)
	if err != nil {
		return err
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		} else {
			ioutil.WriteFile("features.zip", body, os.ModePerm)
			err := unzip("features.zip", ".")
			if err != nil {
				os.Remove("features.zip")
				return err
			} else {
				os.Remove("features.zip")
			}
		}
	}
	return nil

}

// exportTestsExecution export test to provided host
func ExportTestExecution(host, resultsFile, user, password string) error {
	reqUrl := host + "/rest/raven/1.0/import/execution/cucumber"
	client := &http.Client{}

	jsonResult, _ := ioutil.ReadFile(resultsFile)

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(jsonResult))
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")
	// ...
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New("Problem uploading results" + string(body))
	}

	return nil
}

// exportTests export test to selected output directory
func ExecuteTestSet(host string, filter int, outputDirectory, user, password, keys, resultFile string) (err error) {

	DownloadTests(host, filter, outputDirectory, user, password, keys)
	_ = GetPendingCucumberSteps(outputDirectory)
	err1 := ExecuteCucumberTest("json_pretty", resultFile, outputDirectory)
	if err1 != nil {
		return err1
	}
	err1 = ExportTestExecution(host, resultFile, user, password)
	if err1 != nil {
		return err1
	}
	os.Remove(resultFile)
	return nil

}

func ExecuteCucumberTest(format, resultFile, featureDir string) (err error) {
	cmd := "cucumber"
	args := []string{"-x"}
	if format != "" {
		args = append(args, "--format="+format)
	}

	if resultFile != "" {
		args = append(args, "--out="+resultFile)
	}
	args = append(args, featureDir)
	cucumberCommand := exec.Command(cmd, args...)
	cucumberCommand.Stdout = os.Stdout
	cucumberCommand.Stderr = os.Stderr
	if err := cucumberCommand.Run(); err != nil {
		if exists, _ := Exists(resultFile); !exists {
			return err
		}
	}
	return nil
}

func GetPendingCucumberSteps(featureDir string) error {
	cmd := "cucumber"
	args := []string{"--no-color", "-s"}
	args = append(args, featureDir)
	cucumberCommand := exec.Command(cmd, args...)
	cmd = "sed"
	args = []string{"-n", "/You can implement/,/@/p"}
	sedCommand1 := exec.Command(cmd, args...)
	cmd = "sed"
	args = []string{"-e", "/You can implement/c\\"}
	sedCommand2 := exec.Command(cmd, args...)

	sedCommand1.Stdin, _ = cucumberCommand.StdoutPipe()
	sedCommand2.Stdin, _ = sedCommand1.StdoutPipe()
	result, _ := sedCommand2.StdoutPipe()
	_ = cucumberCommand.Start()
	_ = sedCommand1.Start()
	_ = cucumberCommand.Wait()
	_ = sedCommand2.Start()
	_ = sedCommand1.Wait()
	pendingSteps, err := ioutil.ReadAll(result)
	_ = sedCommand2.Wait()

	if err != nil {
		return err
	} else {
		if len(pendingSteps) > 0 {
			current_time := time.Now().Local()
			fileName := "features/step_definitions/pending_" + current_time.Format(time.Stamp) + ".rb"
			fileName = strings.Replace(strings.Replace(fileName, " ", "_", -1), ":", "_", -1)
			StoreResults(fileName, pendingSteps)
		}
	}

	return nil
}
