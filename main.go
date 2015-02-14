package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/octokit/go-octokit/octokit"
)

type MergeResponse struct {
	SHA     string `json:"sha"`
	Merged  bool   `json:"merged"`
	Message string `json:"message"`
}

var (
	verbose             = false
	PullRequestMergeURL = octokit.Hyperlink("repos/{owner}/{repo}/pulls{/number}/merge")
)

func main() {
	verbose = os.Getenv("VERBOSE") != ""

	flag.Parse()

	owner, repo := fetchRepoOwnerAndName()
	number := flag.Arg(0)
	if number == "" {
		fmt.Println("Specify a PR number without the #.")
		os.Exit(1)
	}

	client := octokit.NewClient(octokit.NetrcAuth{})
	url, _ := PullRequestMergeURL.Expand(octokit.M{
		"owner":  owner,
		"repo":   repo,
		"number": number,
	})
	fmt.Printf("Attempting to merge PR #%s on %s/%s...\n", number, owner, repo)

	req, err := client.NewRequest(url.String())
	if err != nil {
		log.Fatal(err)
	}

	var merged MergeResponse
	res, mergeErr := req.Put(map[string]string{}, &merged)

	if mergeErr != nil {
		if verbose {
			fmt.Println("Received an error!", mergeErr)
		}
		if strings.Contains(mergeErr.Error(), "405 - Pull request") {
			fmt.Print("That PR can't be merged. Continue anyway? (y/n) ")
			var answer string
			fmt.Scanln(&answer)
			if answer != "y" {
				return
			}
		} else {
			fmt.Println("Either that's not a pull request or I'm crazy.")
			os.Exit(1)
		}

	}
	log.Println("req", req)
	log.Println("merged", merged)
	log.Println("res", res)

	openEditor()
	commitChangesToHistoryFile(number)
}
