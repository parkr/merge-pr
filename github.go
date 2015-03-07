package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	client *github.Client

	NotMergableError        = errors.New("Not mergable")
	NotFoundError           = errors.New("Not found")
	BranchNotFoundError     = errors.New("Branch not found")
	NonDeletableBranchError = errors.New("Branch cannot be deleted")
)

func init() {
	client = github.NewClient(authenticatedClient())
}

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func netrcPath() string {
	filename := filepath.Join(os.Getenv("HOME"), ".netrc")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		panic(fmt.Sprintf("no such file or directory: %s", filename))
	}
	return filename
}

func accessTokenFromNetrc() string {
	machine, err := netrc.FindMachine(netrcPath(), "api.github.com")
	if err != nil {
		panic(err)
	}
	return machine.Password
}

func authenticatedClient() *http.Client {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: accessTokenFromNetrc()},
	}
	return oauth2.NewClient(oauth2.NoContext, ts)
}

func stringToInt(number string) int {
	intVal, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}
	return intVal
}

func getPullRequest(owner, repo, number string) (*github.PullRequest, error) {
	pr, _, prGetErr := client.PullRequests.Get(owner, repo, stringToInt(number))
	return pr, prGetErr
}

func mergePullRequest(owner, repo, number string) error {
	if verbose {
		log.Printf("Attempting to merge PR #%s on %s/%s...\n", number, owner, repo)
	}

	commitMsg := fmt.Sprintf("Merge pull request %v", number)
	_, _, mergeErr := client.PullRequests.Merge(owner, repo, stringToInt(number), commitMsg)

	if mergeErr != nil {
		if verbose {
			fmt.Println("Received an error!", mergeErr)
		}
		if strings.Contains(mergeErr.Error(), "405 Pull Request is not mergeable") {
			return NotMergableError
		} else {
			if strings.Contains(mergeErr.Error(), "404 Not Found") {
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

	_, deleteBranchErr := client.Git.DeleteRef(owner, repo, branch)

	if strings.Contains(deleteBranchErr.Error(), "422 Reference does not exist") {
		return BranchNotFoundError
	}

	if strings.Contains(deleteBranchErr.Error(), "No media type for this response") {
		return nil
	}

	return deleteBranchErr
}
