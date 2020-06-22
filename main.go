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
	version     = "1.1.2"
	revision    = "dev"
)

func fatalError(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}

func main() {
	flag.BoolVar(&verbose, "v", false, "run verbosely")
	flag.BoolVar(&showVersion, "V", false, "print version and exit")
	flag.Parse()

	if showVersion {
		fmt.Printf("merge-pr %s (%s)\n", version, revision)
		os.Exit(0)
	}

	number := flag.Arg(0)
	if number == "" {
		fatalError("Specify a PR number without the #.")
	}

	if verbose {
		log.Println("Determining if your local branch is cool.")
	}
	err := isAcceptableCurrentBranch()
	if err != nil {
		fatalError(err.Error())
	}

	initializeGitHubClient()

	if verbose {
		log.Println("Fetching owner & repo from your git remotes")
	}
	owner, repo := fetchRepoOwnerAndName()
	if owner == "" || repo == "" {
		fatalError("You don't have an 'origin' remote. Failing.")
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
			fatalError(err.Error())
		}
	}

	if err := gitPull(); err != nil {
		log.Fatal(err)
	}
	openEditor()
	commitChangesToHistoryFile(number)
	gitPush()
}
