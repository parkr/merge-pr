package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var (
	client *github.Client

	NotMergableError        = errors.New("Not mergable")
	BranchNotFoundError     = errors.New("Branch not found")
	NonDeletableBranchError = errors.New("Branch cannot be deleted")
	PullReqNotFoundError    = errors.New("Pull request not found")
)

func initializeGitHubClient() {
	client = github.NewClient(newClient())
}

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func hubConfigPath() string {
	filename := filepath.Join(os.Getenv("HOME"), ".config", "hub")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if verbose {
			fmt.Printf("no such file or directory: %s", filename)
		}
		return ""
	}
	return filename
}

func accessTokenFromHubConfig() string {
	f, err := os.Open(hubConfigPath())
	if err != nil {
		return ""
	}

	config := struct {
		GitHub []struct {
			OauthToken string `yaml:"oauth_token"`
		} `yaml:"github.com"`
	}{}
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Printf("couldn't decode hub config: %+v", err)
		return ""
	}
	if len(config.GitHub) == 0 {
		return ""
	}
	return config.GitHub[0].OauthToken
}

func netrcPath() string {
	filename := filepath.Join(os.Getenv("HOME"), ".netrc")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if verbose {
			fmt.Printf("no such file or directory: %s", filename)
		}
		return ""
	}
	return filename
}

func accessTokenFromNetrc() string {
	if netrcPath() == "" {
		return ""
	}

	machine, err := netrc.FindMachine(netrcPath(), "api.github.com")
	if err != nil {
		panic(err)
	}
	if machine == nil {
		return ""
	}
	return machine.Password
}

func newClient() *http.Client {
	if accessToken := accessTokenFromNetrc(); accessToken != "" {
		return oauth2.NewClient(oauth2.NoContext, &tokenSource{
			&oauth2.Token{AccessToken: accessToken},
		})
	}
	if accessToken := accessTokenFromHubConfig(); accessToken != "" {
		return oauth2.NewClient(oauth2.NoContext, &tokenSource{
			&oauth2.Token{AccessToken: accessToken},
		})
	}
	return http.DefaultClient
}

func stringToInt(number string) int {
	intVal, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}
	return intVal
}

func getPullRequest(owner, repo, number string) (*github.PullRequest, error) {
	pr, res, prGetErr := client.PullRequests.Get(context.Background(), owner, repo, stringToInt(number))
	if prGetErr != nil {
		switch res.StatusCode {
		case http.StatusNotFound:
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
	_, res, mergeErr := client.PullRequests.Merge(
		context.Background(),
		owner,
		repo,
		stringToInt(number),
		commitMsg,
		&github.PullRequestOptions{},
	)

	if mergeErr != nil {
		if verbose {
			fmt.Println("Received an error!", mergeErr)
		}
		// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#merge-a-pull-request--status-codes
		switch res.StatusCode {
		case http.StatusMethodNotAllowed:
			return NotMergableError
		case http.StatusNotFound:
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
	case "main", "master", "gh-pages", "dev", "staging":
		return NonDeletableBranchError
	}

	ref := fmt.Sprintf("heads/%s", branch)

	res, deleteBranchErr := client.Git.DeleteRef(context.Background(), owner, repo, ref)

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
