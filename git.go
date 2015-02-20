package main

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var (
	GitRemoteRegexp = regexp.MustCompile("(https|git)(@|://)github.com(:|/)([a-zA-Z0-9-_]+)/([a-zA-Z0-9-_]+).git")
)

func fetchRepoOwnerAndName() (string, string) {
	return extractOwnerAndNameFromRemote(gitOriginRemote())
}

func extractOwnerAndNameFromRemote(url string) (string, string) {
	matches := GitRemoteRegexp.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", ""
	}
	return matches[len(matches)-2], matches[len(matches)-1]
}

func gitOriginRemote() string {
	out, err := exec.Command("git", "config", "remote.origin.url").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(string(out), "\n")
}

func gitPull() {
	shellExec("git", "pull", "--rebase")
}

func gitPush() {
	shellExec("git", "push")
}

func commitChangesToHistoryFile(pr string) {
	shellExec("git", "add", "History.markdown")
	shellExec(
		"git",
		"commit",
		"-m",
		"Update history to reflect merge of #"+pr+" [ci skip]",
	)
}
