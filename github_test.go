package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
)

var isCI = os.Getenv("CI") == "1"

func initializeTestClient(handler http.Handler) *httptest.Server {
	server := httptest.NewServer(handler)
	u, _ := url.Parse(server.URL)
	u.Path = "/v3/"

	client = github.NewClient(newClient())
	client.BaseURL = u

	return server
}

func TestMergePullRequest_WithIssue(t *testing.T) {
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v3/repos/parkr/merge-pr/pulls/2/merge", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		w.WriteHeader(http.StatusNotFound)
	})
	server := initializeTestClient(mux)
	defer server.Close()

	err := mergePullRequest("parkr", "merge-pr", "2")

	assert.EqualError(t, err, "Pull request not found")
}

func TestMergePullRequest_WithAlreadyMergedPR(t *testing.T) {
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v3/repos/parkr/merge-pr/pulls/1/merge", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	server := initializeTestClient(mux)
	defer server.Close()

	err := mergePullRequest("parkr", "merge-pr", "1")

	assert.EqualError(t, err, "Not mergable")
}

func TestGetPullRequest(t *testing.T) {
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v3/repos/parkr/merge-pr/pulls/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Write([]byte(`{"head":{"ref":"do-it-all","user":{"login":"parkr"}}}`))
	})
	server := initializeTestClient(mux)
	defer server.Close()

	pr, err := getPullRequest("parkr", "merge-pr", "1")

	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, "do-it-all", *pr.Head.Ref)
	assert.Equal(t, "parkr", *pr.Head.User.Login)
}

func TestDeleteBranchForPR_BranchNotFound(t *testing.T) {
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/repos/parkr/merge-pr/pulls/1" {
			assert.Equal(t, "GET", r.Method)
			w.Write([]byte(`{"head":{"ref":"do-it-all","user":{"login":"parkr"}}}`))
			return
		}
		if r.URL.Path == "/v3/repos/parkr/merge-pr/git/refs/heads/do-it-all" {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		assert.Equal(t, "not expected", r.URL.Path)
		assert.Equal(t, "not expected", r.Method)
		w.Write([]byte(`{}`))
	})
	server := initializeTestClient(mux)
	defer server.Close()

	err := deleteBranchForPullRequest("parkr", "merge-pr", "1")

	assert.EqualError(t, err, "Branch not found")
}

func TestDeleteBranchForPR_ForNonPR(t *testing.T) {
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/repos/parkr/merge-pr/pulls/2" {
			assert.Equal(t, "GET", r.Method)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"head":{"ref":"do-it-all","user":{"login":"parkr"}}}`))
			return
		}
		assert.Equal(t, "expected", r.URL.Path)
		assert.Equal(t, "expected", r.Method)
		w.Write([]byte(`{}`))
	})
	server := initializeTestClient(mux)
	defer server.Close()

	err := deleteBranchForPullRequest("parkr", "merge-pr", "2")

	assert.EqualError(t, err, "Pull request not found")
}

func TestDeleteBranch(t *testing.T) {
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/repos/parkr/merge-pr/git/refs/heads/do-it-all" {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	})
	server := initializeTestClient(mux)
	defer server.Close()

	err := deleteBranch("parkr", "merge-pr", "do-it-all")

	assert.EqualError(t, err, "Branch not found")
}

func TestDeleteBranch_WithProtectedBranch(t *testing.T) {
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "expected", r.URL.Path)
		assert.Equal(t, "expected", r.Method)
		w.Write([]byte(`{}`))
	})
	server := initializeTestClient(mux)
	defer server.Close()

	err := deleteBranch("parkr", "merge-pr", "gh-pages")

	assert.EqualError(t, err, "Branch cannot be deleted")
}
