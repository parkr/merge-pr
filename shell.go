package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func commandFromArgs(args ...string) *exec.Cmd {
	return exec.Command(args[0], args[1:len(args)]...)
}

func shellExec(args ...string) error {
	if verbose {
		log.Println(args)
	}
	cmd := commandFromArgs(args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func shellOutput(args ...string) string {
	out, err := commandFromArgs(args...).Output()
	if err != nil {
		fatalError(err.Error())
	}
	return strings.TrimRight(string(out), "\n")
}
