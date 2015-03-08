package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	client *github.Client

	NotMergableError        = errors.New("Not mergable")
	BranchNotFoundError     = errors.New("Branch not found")
	NonDeletableBranchError = errors.New("Branch cannot be deleted")
	PullReqNotFoundError    = errors.New("Pull request not found")
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
	pr, res, prGetErr := client.PullRequests.Get(owner, repo, stringToInt(number))
	if prGetErr != nil {
		switch res.StatusCode {
		case 404:
			return nil, PullReqNotFoundError
		default:
			return nil, prGetErr
		}
	}
	return pr, prGetErr
}

func mergePullRequest(owner, repo, number string) error {
	if verbose {
		log.Printf("Attempting to merge PR #%s on %s/%s...\n", number, owner, repo)
	}

	commitMsg := fmt.Sprintf("Merge pull request %v", number)
	_, res, mergeErr := client.PullRequests.Merge(owner, repo, stringToInt(number), commitMsg)

	if mergeErr != nil {
		if verbose {
			fmt.Println("Received an error!", mergeErr)
		}
		switch res.StatusCode {
		case 405:
			return NotMergableError
		case 404:
			return PullReqNotFoundError
		default:
			return mergeErr
		}
	}

	return nil
}

func deleteBranchForPullRequest(owner, repo, number string) error {
	pr, prGetErr := getPullRequest(owner, repo, number)
	if prGetErr != nil {
		return prGetErr
	}

	if *pr.Head.User.Login == owner && *pr.Head.Ref != "" {
		if verbose {
			log.Println("Deleting the branch.")
		}
		err := deleteBranch(owner, repo, *pr.Head.Ref)
		if err != nil {
			return err
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

	res, deleteBranchErr := client.Git.DeleteRef(owner, repo, branch)

	if deleteBranchErr != nil {
		switch res.StatusCode {
		case 422:
			return BranchNotFoundError
		default:
			return deleteBranchErr
		}
	}

	return nil
}
