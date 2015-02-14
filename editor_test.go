package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHistoryFile(t *testing.T) {
	assert.True(t, isHistoryFile("History.markdown"))
	assert.True(t, isHistoryFile("HISTORY.markdown"))
	assert.True(t, isHistoryFile("History.mkd"))
	assert.True(t, isHistoryFile("History.mkdn"))
	assert.True(t, isHistoryFile("History.md"))
	assert.True(t, isHistoryFile("History.MD"))
	assert.True(t, isHistoryFile("HISTORY.MD"))
	assert.True(t, isHistoryFile("History.MKDN"))
	assert.True(t, isHistoryFile("Changelog.mkdn"))
	assert.True(t, isHistoryFile("CHANGELOG.markdown"))
	assert.True(t, isHistoryFile("CHANGELOG.MD"))
	assert.True(t, isHistoryFile("Changelog.mkd"))
}

func TestFindHistoryFile(t *testing.T) {
	assert.Equal(t, "History.markdown", historyFile())
}
