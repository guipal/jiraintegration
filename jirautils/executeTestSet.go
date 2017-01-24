// Package jiraUtils contains utility functions for working with strings.
package jirautils

import (
	"os"

	"os/exec"

	"fmt"
)

// exportTests export test to selected output directory
func ExecuteTestSet(host string, filter int, outputDirectory, user, password, keys, resultFile string) (err error) {

	ExportTests(host, filter, outputDirectory, user, password, keys)
	err1 := ExecuteCucumberTest(resultFile)
	if err1 != nil {
		return err1
	}
	err1 = ImportTestsExecution(host, resultFile, user, password)
	if err1 != nil {
		return err1
	}
	os.Remove(resultFile)
	return nil

}

func ExecuteCucumberTest(resultFile string) (err error) {
	fmt.Println("Executing cucumber test")
	cmd := "cucumber"
	args := []string{"-x", "--format=json_pretty", "--out=" + resultFile}
	cucumberCommand := exec.Command(cmd, args...)
	cucumberCommand.Stderr = os.Stderr
	if err := cucumberCommand.Run(); err != nil {
		if exists, _ := exists(resultFile); !exists {
			return err
		}
	}
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
