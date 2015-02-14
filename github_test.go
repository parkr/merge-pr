package main

import (
	"testing"

	"github.com/octokit/go-octokit/octokit"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	req, err := newRequest("parkr", "merge-pr", "1234")
	assert.NoError(t, err)
	assert.IsType(t, &octokit.Request{}, req)
}

func TestMergePullRequestWithIssue(t *testing.T) {
	err := mergePullRequest("parkr", "merge-pr", "2")
	assert.EqualError(t, err, "Not found")
}

func TestMergePullRequestWithAlreadyMergedPR(t *testing.T) {
	err := mergePullRequest("parkr", "merge-pr", "1")
	assert.EqualError(t, err, "Not mergable")
}
