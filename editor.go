package main

import (
	"os"
	"os/exec"
)

func openEditor() {
	cmd := exec.Command(os.Getenv("EDITOR"), "History.markdown")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
