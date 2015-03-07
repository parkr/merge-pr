package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var verbose = false

func init() {
	flag.BoolVar(&verbose, "v", false, "run verbosely")
}

func main() {
	flag.Parse()

	number := flag.Arg(0)
	if number == "" {
		fmt.Println("Specify a PR number without the #.")
		os.Exit(1)
	}
	if verbose {
		log.Println("Fetching owner & repo from your git remotes")
	}
	owner, repo := fetchRepoOwnerAndName()
	if owner == "" || repo == "" {
		fmt.Println("You don't have an 'origin' remote. Failing.")
		os.Exit(1)
	}

	if verbose {
		log.Println("Attempting to merge the PR.")
	}
	err := mergePullRequest(owner, repo, number)
	if err != nil {
		if err == NotMergableError {
			fmt.Print("That PR can't be merged. Continue anyway? (y/n) ")
			var answer string
			fmt.Scanln(&answer)
			if answer != "y" {
				os.Exit(1)
			}
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if verbose {
		log.Println("Grabbing the PR's data.")
	}
	pr, err := getPullRequest(owner, repo, number)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *pr.Head.User.Login == owner && *pr.Head.Ref != "" {
		if verbose {
			log.Println("Deleting the branch.")
		}
		err := deleteBranch(owner, repo, *pr.Head.Ref)
		if err != nil {
			fmt.Println(err)
		}
	}

	gitPull()
	openEditor()
	commitChangesToHistoryFile(number)
	gitPush()
}
