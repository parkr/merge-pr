package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	verbose     bool
	showVersion bool
	version     = "1.1.1"
)

func init() {
	flag.BoolVar(&verbose, "v", false, "run verbosely")
	flag.BoolVar(&showVersion, "V", false, "print version and exit")
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("merge-pr %v\n", version)
		os.Exit(0)
	}

	number := flag.Arg(0)
	if number == "" {
		fmt.Println("Specify a PR number without the #.")
		os.Exit(1)
	}

	if verbose {
		log.Println("Determining if your local branch is cool.")
	}
	err := isAcceptableCurrentBranch()
	if err != nil {
		log.Fatal(err)
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
	err = mergePullRequest(owner, repo, number)
	if err == nil {
		if verbose {
			log.Println("Deleting branch for PR.")
		}
		err = deleteBranchForPullRequest(owner, repo, number)
		if err != nil {
			fmt.Println("Error deleting the branch:", err)
		}
	} else {
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

	if err := gitPull(); err != nil {
		log.Fatal(err)
	}
	openEditor()
	commitChangesToHistoryFile(number)
	gitPush()
}
