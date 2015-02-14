package main

import (
	"os"
	"testing"

	"github.com/octokit/go-octokit/octokit"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	req, err := newMergeRequest("parkr", "merge-pr", "1234")
	assert.NoError(t, err)
	assert.IsType(t, &octokit.Request{}, req)
}

func TestMergePullRequestWithIssue(t *testing.T) {
	if os.Getenv("USER") != "travis" {
		err := mergePullRequest("parkr", "merge-pr", "2")
		assert.EqualError(t, err, "Not found")
	}
}

func TestMergePullRequestWithAlreadyMergedPR(t *testing.T) {
	if os.Getenv("USER") != "travis" {
		err := mergePullRequest("parkr", "merge-pr", "1")
		assert.EqualError(t, err, "Not mergable")
	}
}

func TestGetPullRequest(t *testing.T) {
	if os.Getenv("USER") != "travis" {
		pr, err := getPullRequest("parkr", "merge-pr", "1")
		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, "do-it-all", pr.Head.Ref)
		assert.Equal(t, "parkr", pr.Head.User.Login)
	}
}

func TestDeleteBranch(t *testing.T) {
	if os.Getenv("USER") != "travis" {
		err := deleteBranch("parkr", "merge-pr", "do-it-all")
		assert.EqualError(t, err, "Branch not found")
	}
}

func TestDeleteBranchWithProtectedBranch(t *testing.T) {
	if os.Getenv("USER") != "travis" {
		err := deleteBranch("parkr", "merge-pr", "gh-pages")
		assert.EqualError(t, err, "Branch cannot be deleted")
	}
}
