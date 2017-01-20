// Package jiraUtils contains utility functions for working with strings.
package jirautils

import (
	"os"

	"net/http"

	"fmt"

	"strconv"

	"io"

	"io/ioutil"

	"archive/zip"

	"path/filepath"
)

// exportTests export test to selected output directory
func ExportTests(host string, filter int, outputDirectory string) {
	os.MkdirAll(outputDirectory, os.ModePerm)
	os.Chdir(outputDirectory)

	client := &http.Client{}
	req, err := http.NewRequest("GET", host+"/jira/rest/raven/1.0/export/test?filter="+strconv.Itoa(filter), nil)
	// ...
	req.Header.Add("Authorization", "Basic bHV6aGVvbGc6QWRpLUFxYV8yMDE2")
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	} else {
		fmt.Println("Response from server:", resp)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		} else {
			ioutil.WriteFile("features.zip", body, os.ModePerm)
			err := Unzip("features.zip", ".")
			if err != nil {
				os.Remove("features.zip")

			} else {
				os.Remove("features.zip")
			}
		}
	}

}
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
