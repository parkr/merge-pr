package main

import (
	"flag"
	"fmt"
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
	owner, repo := fetchRepoOwnerAndName()
	if owner == "" || repo == "" {
		fmt.Println("You don't have an 'origin' remote. Failing.")
		os.Exit(1)
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

	pr, err := getPullRequest(owner, repo, number)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if pr.Head.User.Login == owner && pr.Head.Ref != "" {
		err := deleteBranch(owner, repo, pr.Head.Ref)
		if err != nil {
			fmt.Println(err)
		}
	}

	gitPull()
	openEditor()
	commitChangesToHistoryFile(number)
}
