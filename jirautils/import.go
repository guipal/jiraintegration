// Package jiraUtils contains utility functions for working with strings.
package jirautils

import (
	"net/http"

	"io/ioutil"

	"bytes"

	"errors"
)

// exportTests export test to selected output directory
func ImportTestsExecution(host, resultsFile, user, password string) (err error) {
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
		return errors.New("Problem uploading results")
	}

	return nil
}
