package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/octokit/go-octokit/octokit"
)

var (
	mergePullRequestUrl = octokit.Hyperlink("repos/{owner}/{repo}/pulls{/number}/merge")

	NotMergableError = errors.New("Not mergable")
	NotFoundError    = errors.New("Not found")
)

func newRequest(owner, repo, number string) (*octokit.Request, error) {
	client := octokit.NewClient(octokit.NetrcAuth{})
	url, _ := mergePullRequestUrl.Expand(octokit.M{
		"owner":  owner,
		"repo":   repo,
		"number": number,
	})

	return client.NewRequest(url.String())
}

func mergePullRequest(owner, repo, number string) error {
	if verbose {
		log.Printf("Attempting to merge PR #%s on %s/%s...\n", number, owner, repo)
	}

	req, err := newRequest(owner, repo, number)
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
