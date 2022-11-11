package main

import (
	"fmt"
	"regexp"
)

var (
	GitRemoteRegexp    = regexp.MustCompile("(https|git)(@|://)github\\.com(:|/)([a-zA-Z0-9-_]+)/([a-zA-Z0-9-_]+)(?:\\.git)?")
	acceptableBranches = []string{"master", "main", "staging", "dev"}
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func currentBranch() string {
	return shellOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
}

func isAcceptableCurrentBranch() error {
	currBranch := currentBranch()
	if !contains(acceptableBranches, currBranch) {
		return fmt.Errorf("Unacceptable local branch: %s", currBranch)
	}
	return nil
}

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
	return shellOutput("git", "config", "remote.origin.url")
}

func gitPull() error {
	return shellExec("git", "pull", "--rebase")
}

func gitPush() {
	shellExec("git", "push")
}

func commitChangesToHistoryFile(pr string) {
	shellExec("git", "add", historyFile())
	shellExec(
		"git",
		"commit",
		"-m",
		"Update history to reflect merge of #"+pr,
		"-m",
		"[ci skip]",
	)
}
