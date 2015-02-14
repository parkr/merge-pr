package main

import (
	"log"

	"github.com/octokit/go-octokit/octokit"
)

var (
	PullRequestMergeURL = octokit.Hyperlink("repos/{owner}/{repo}/pulls{/number}/merge")
)

func main() {
	client := octokit.NewClient(octokit.NetrcAuth{})
	url, _ := PullRequestMergeURL.Expand(octokit.M{
		"owner":  "jekyll",
		"repo":   "jekyll",
		"number": 3455,
	})
    log.Println(url)

	sawyerReq, err := client.Client.NewRequest(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	req = &Request{client: client, Request: sawyerReq}
    log.Println(req)
}
