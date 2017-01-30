// Package jiraUtils contains utility functions for working with strings.

package jirautils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func unzip(src, dest string) error {
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
func StoreResults(path string, result []byte) error {
	f, _ := os.Create(path)
	defer f.Close()
	_, err := f.Write(result)
	f.Sync()
	if err != nil {
		return err
	}
	return nil

}

func GetRepoName(repo string) string {
	c1 := exec.Command("echo", repo)
	c2 := exec.Command("sed", "s/.*[:/]\\([^:/]*\\)\\.git$/\\1/")

	c2.Stdin, _ = c1.StdoutPipe()
	result, _ := c2.StdoutPipe()
	_ = c1.Start()
	_ = c2.Start()
	_ = c1.Wait()
	repoName, _ := ioutil.ReadAll(result)
	_ = c2.Wait()

	return strings.TrimSpace(string(repoName))

}

func cloneRepo(repoName string, repo string) {

	cmd := "git"
	args := []string{"clone", repo}
	fmt.Println("Cloning repo: ", repoName)
	cloneCommand := exec.Command(cmd, args...)
	cloneCommand.Stdin = os.Stdin
	cloneCommand.Stdout = os.Stdout
	cloneCommand.Stderr = os.Stderr
	if err := cloneCommand.Run(); err != nil {
		os.Exit(1)
	}
	fmt.Println("Successfully cloned", repoName)

}
