# merge-pr

Merge your GitHub pull requests from the command line.

[![Build Status](https://travis-ci.org/parkr/merge-pr.svg?branch=master)](https://travis-ci.org/parkr/merge-pr)

## Motivation

Merging pull requests in the browser is nice, sure, but you lose a lot of
clarity into what changed & when at a higher level. When did you add that
feature? Oh, you're releasing a new patch? What changes did you make? No, I
won't read through your commit history.

This tool aims to make it easy to merge PR's and add a line to the
CHANGELOG (`History.markdown` by default). All with one command.

Here's more on [why you should keep a changelog.](http://keepachangelog.com/)

## Installation

You need [Go](https://golang.org) and you need your `$GOPATH` set &
`$GOPATH/bin` in your `$PATH`. Then:

```bash
$ go get github.com/parkr/merge-pr
```

Throw your credentials in `$HOME/.netrc`, like this:

```text
machine api.github.com
  login yourusername
  password mypersonalaccesstokenforgithub
```

Grab a personal access token on the [GitHub Applications Setting
page](https://github.com/settings/applications).

## Usage

```bash
$ cd my-project
$ merge-pr 7
```

It uses the origin remote to discern which repo to make the API requests
against, so ensure your `origin` is pointed to the repository that you
want to merge the pull request into.

This will go to GitHub, merge the PR, delete the branch if it's on the same
repo, will pull down those changes, open up your editor (`$EDITOR`), then
commit that change.

## Contributing

To get setup, clone the repo and run `script/bootstrap`. Make your edits.
Add tests where you can. Run tests with `script/test`, use `script/cibuild`
to run the tests and build the binary if the tests are successful.

Once you're happy with your change, submit a PR. If I like it, I'll use
this tool to merge it!

## Versioning

We adhere to SemVer where applicable. To see the version of your copy of
`merge-pr`, run `merge-pr -V`.

To release a new version, change the version in `main.go`, and run
`script/release`.

## Credits / License

MIT License, copyright Parker Moore. Details in the `LICENSE` file.
