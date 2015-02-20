package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellExec(t *testing.T) {
	assert.NoError(t, shellExec("echo"))
}

func TestShellExecForMultiArgs(t *testing.T) {
	assert.NoError(t, shellExec("sh", "-c", "test -f History.markdown"))
}

func TestShellExecFailingCommand(t *testing.T) {
	err := shellExec("sh", "-c", "test", "-d", "History.markdown")
	assert.EqualError(t, err, "exit status 1")
}
