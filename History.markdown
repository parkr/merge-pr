## HEAD

  * git: don't requre a `.git` suffix for the git remote URL (#26)

## 1.1.2 / 2015-05-12

  * git: commit computed history file instead of hard-coded history.markdown (#24)

## 1.1.1 / 2015-05-07

  * git: put `[ci skip]` in the commit msg body instead of the summary (#22)

## 1.1.0 / 2015-03-11

  * git: if 'git pull' fails, do not continue (#21)
  * fix formatting for version printing (#20)
  * git: only merge if the current branch is master, staging, or dev (#19)
  * github: delete the branch properly (#18)
  * Replace go-octokit with go-github. (#17)

## 1.0.0 / 2015-02-19

  * git: use 'git config' to get origin URL (#12)
  * Add more printed info with the `-v` flag (#9)
  * Shell out to the OS with a common interface (#8)
  * Push once the merge commit is added. (#7)
  * Make the fetching of the history file more flexible. (#6)
  * When running the commit command, just use Stdout. (#5)
  * Delete the branch once it's been merged. (#4)
  * Configure to work with Travis (#3)
  * Birthday! (#1)
