package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	gitRemoteSSH    = "git@github.com:parkr/merge-pr.git"
	gitRemoteHTTPS  = "https://github.com/parkr/merge-pr.git"
	gitRemoteGit    = "git://github.com/parkr/merge-pr.git"
	gitRemoteRegexp = regexp.MustCompile("git(@|://)github.com(:|/)parkr/merge-pr.git")
)

func TestOriginRemote(t *testing.T) {
	origin := gitOriginRemote()
	assert.Regexp(t, gitRemoteRegexp, origin)
}

func TestExtractOwnerAndRepoWithSSHUrl(t *testing.T) {
	owner, repo := extractOwnerAndNameFromRemote(gitRemoteSSH)
	assert.Equal(t, "parkr", owner)
	assert.Equal(t, "merge-pr", repo)
}

func TestExtractOwnerAndRepoWithGitUrl(t *testing.T) {
	owner, repo := extractOwnerAndNameFromRemote(gitRemoteGit)
	assert.Equal(t, "parkr", owner)
	assert.Equal(t, "merge-pr", repo)
}

func TestExtractOwnerAndRepoWithHTTPSUrl(t *testing.T) {
	owner, repo := extractOwnerAndNameFromRemote(gitRemoteHTTPS)
	assert.Equal(t, "parkr", owner)
	assert.Equal(t, "merge-pr", repo)
}

func TestExtractOwnerAndRepoWithBadURL(t *testing.T) {
	owner, repo := extractOwnerAndNameFromRemote("git@github.com:L!!!!RS/mars.git")
	assert.Equal(t, "", owner)
	assert.Equal(t, "", repo)
}

func TestFetchRepoOwnerAndName(t *testing.T) {
	owner, repo := fetchRepoOwnerAndName()
	assert.Contains(t, owner, "parkr")
	assert.Contains(t, repo, "merge-pr")
}
