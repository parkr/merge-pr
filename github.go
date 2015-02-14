package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/octokit/go-octokit/octokit"
)

var (
	client *octokit.Client

	mergePullRequestUrl = octokit.Hyperlink("repos/{owner}/{repo}/pulls{/number}/merge")
	deleteBranchUrl     = octokit.Hyperlink("repos/{owner}/{repo}/git/refs/heads/{branch}")

	NotMergableError        = errors.New("Not mergable")
	NotFoundError           = errors.New("Not found")
	BranchNotFoundError     = errors.New("Branch not found")
	NonDeletableBranchError = errors.New("Branch cannot be deleted")
)

func init() {
	client = octokit.NewClient(octokit.NetrcAuth{})
}

func newMergeRequest(owner, repo, number string) (*octokit.Request, error) {
	url, _ := mergePullRequestUrl.Expand(octokit.M{
		"owner":  owner,
		"repo":   repo,
		"number": number,
	})

	return client.NewRequest(url.String())
}

func getPullRequest(owner, repo, number string) (*octokit.PullRequest, error) {
	url, _ := octokit.PullRequestsURL.Expand(octokit.M{
		"owner":  owner,
		"repo":   repo,
		"number": number,
	})
	req, err := client.NewRequest(url.String())
	if err != nil {
		return nil, err
	}

	var pullRequest octokit.PullRequest
	_, prGetErr := req.Get(&pullRequest)

	return &pullRequest, prGetErr
}

func mergePullRequest(owner, repo, number string) error {
	if verbose {
		log.Printf("Attempting to merge PR #%s on %s/%s...\n", number, owner, repo)
	}

	req, err := newMergeRequest(owner, repo, number)
	if err != nil {
		return err
	}

	var merged map[string]interface{}
	_, mergeErr := req.Put(map[string]string{}, &merged)

	if mergeErr != nil {
		if verbose {
			fmt.Println("Received an error!", mergeErr)
		}
		if strings.Contains(mergeErr.Error(), "405 - Pull Request is not mergeable") {
			return NotMergableError
		} else {
			if strings.Contains(mergeErr.Error(), "404 - Not Found") {
				return NotFoundError
			} else {
				return mergeErr
			}
		}

	}

	return nil
}

func deleteBranch(owner, repo, branch string) error {
	switch branch {
	case "master":
		fallthrough
	case "gh-pages":
		fallthrough
	case "dev":
		fallthrough
	case "staging":
		return NonDeletableBranchError
	}

	url, _ := deleteBranchUrl.Expand(octokit.M{
		"owner":  owner,
		"repo":   repo,
		"branch": branch,
	})
	req, err := client.NewRequest(url.String())
	if err != nil {
		return err
	}

	var deleted map[string]interface{}
	_, deleteBranchErr := req.Delete(&deleted)

	if strings.Contains(deleteBranchErr.Error(), "422 - Reference does not exist") {
		return BranchNotFoundError
	}

	return deleteBranchErr
}
