package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

var historyFilenameRegexp = regexp.MustCompile("(?i:(History|Changelog).m(ar)?k?d(own)?)")

func openEditor() {
	cmd := exec.Command(os.Getenv("EDITOR"), historyFile())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func historyFile() string {
	infos, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Problem finding your history file.")
		os.Exit(1)
	}
	for _, info := range infos {
		if isHistoryFile(info.Name()) {
			return info.Name()
		}
	}
	return "History.markdown"
}

func isHistoryFile(filename string) bool {
	return historyFilenameRegexp.FindString(filename) != ""
}
