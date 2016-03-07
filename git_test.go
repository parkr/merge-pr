package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	gitRemoteSSH   = "git@github.com:parkr/merge-pr.git"
	gitRemoteHTTPS = "https://github.com/parkr/merge-pr.git"
	gitRemoteGit   = "git://github.com/parkr/merge-pr.git"
)

func TestCurrentBranch(t *testing.T) {
	branch := currentBranch()
	assert.NotEmpty(t, branch)
}

func TestIsAcceptableCurrentBranch(t *testing.T) {
	branch := currentBranch()
	err := isAcceptableCurrentBranch()
	switch branch {
	case "master":
		fallthrough
	case "staging":
		fallthrough
	case "dev":
		assert.NoError(t, err)
	default:
		assert.EqualError(t, err, fmt.Sprintf("Unacceptable local branch: %s", branch))
	}
}

func TestOriginRemote(t *testing.T) {
	assert.Regexp(t, GitRemoteRegexp, gitRemoteGit)
	assert.Regexp(t, GitRemoteRegexp, gitRemoteSSH)

	origin := gitOriginRemote()
	assert.Regexp(t, GitRemoteRegexp, origin)
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
