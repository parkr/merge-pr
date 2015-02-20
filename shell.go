package main

import (
	"log"
	"os"
	"os/exec"
)

func shellExec(args ...string) error {
	if verbose {
		log.Println(args)
	}
	cmd := exec.Command(args[0], args[1:len(args)]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
