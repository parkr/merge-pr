package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	gitRemoteSSH   = "origin\tgit@github.com:parkr/merge-pr.git\t(push)"
	gitRemoteHTTPS = "origin\thttps://github.com/parkr/merge-pr.git\t(push)"
	gitRemoteGit   = "origin\tgit://github.com/parkr/merge-pr.git\t(fetch)"
)

func TestOriginRemoteWithOneRemote(t *testing.T) {
	remotes := []string{"origin\tgit@github.com:parkr/merge-pr.git\t(fetch)"}
	origin := OriginRemote(remotes)
	assert.Equal(t, "git@github.com:parkr/merge-pr.git", origin)
}

func TestOriginRemoteWithLotsOfRemotes(t *testing.T) {
	remotes := []string{
		gitRemoteGit,
		gitRemoteSSH,
	}
	origin := OriginRemote(remotes)
	assert.Equal(t, "git://github.com/parkr/merge-pr.git", origin)
}

func TestExtractRemoteWithSSHUrl(t *testing.T) {
	url := ExtractRemote(gitRemoteSSH)
	assert.Equal(t, "git@github.com:parkr/merge-pr.git", url)
}

func TestExtractRemoteWithGitUrl(t *testing.T) {
	url := ExtractRemote(gitRemoteGit)
	assert.Equal(t, "git://github.com/parkr/merge-pr.git", url)
}

func TestExtractOwnerWithHTTPSUrl(t *testing.T) {
	url := ExtractRemote(gitRemoteHTTPS)
	assert.Equal(t, "https://github.com/parkr/merge-pr.git", url)
}

func TestExtractOwnerAndRepoWithSSHUrl(t *testing.T) {
	owner, repo := ExtractOwnerAndRepo(gitRemoteSSH)
	assert.Equal(t, "parkr", owner)
	assert.Equal(t, "merge-pr", repo)
}

func TestExtractOwnerAndRepoWithGitUrl(t *testing.T) {
	owner, repo := ExtractOwnerAndRepo(gitRemoteGit)
	assert.Equal(t, "parkr", owner)
	assert.Equal(t, "merge-pr", repo)
}

func TestExtractOwnerAndRepoWithHTTPSUrl(t *testing.T) {
	owner, repo := ExtractOwnerAndRepo(gitRemoteHTTPS)
	assert.Equal(t, "parkr", owner)
	assert.Equal(t, "merge-pr", repo)
}

func TestExtractOwnerAndRepoWithBadURL(t *testing.T) {
	owner, repo := ExtractOwnerAndRepo("git@github.com:L!!!!RS/mars.git")
	assert.Equal(t, "", owner)
	assert.Equal(t, "", repo)
}
