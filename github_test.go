package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var isCI = os.Getenv("CI") == "1"

func TestMergePullRequestWithIssue(t *testing.T) {
	if !isCI {
		err := mergePullRequest("parkr", "merge-pr", "2")
		assert.EqualError(t, err, "Pull request not found")
	}
}

func TestMergePullRequestWithAlreadyMergedPR(t *testing.T) {
	if !isCI {
		err := mergePullRequest("parkr", "merge-pr", "1")
		assert.EqualError(t, err, "Not mergable")
	}
}

func TestGetPullRequest(t *testing.T) {
	if !isCI {
		pr, err := getPullRequest("parkr", "merge-pr", "1")
		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, "do-it-all", *pr.Head.Ref)
		assert.Equal(t, "parkr", *pr.Head.User.Login)
	}
}

func TestDeleteBranchForPR(t *testing.T) {
	if !isCI {
		err := deleteBranchForPullRequest("parkr", "merge-pr", "1")
		assert.EqualError(t, err, "Branch not found")
	}
}

func TestDeleteBranchForPRForNonPR(t *testing.T) {
	if !isCI {
		err := deleteBranchForPullRequest("parkr", "merge-pr", "2")
		assert.EqualError(t, err, "Pull request not found")
	}
}

func TestDeleteBranch(t *testing.T) {
	if !isCI {
		err := deleteBranch("parkr", "merge-pr", "do-it-all")
		assert.EqualError(t, err, "Branch not found")
	}
}

func TestDeleteBranchWithProtectedBranch(t *testing.T) {
	if !isCI {
		err := deleteBranch("parkr", "merge-pr", "gh-pages")
		assert.EqualError(t, err, "Branch cannot be deleted")
	}
}
