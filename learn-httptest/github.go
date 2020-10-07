package learn_httptest

import (
	"context"
	"net/http"
	"github.com/google/go-github/v32/github"
)

const owner = "an-owner"
const repo = "a-repo"

type Github struct {
	Client *github.Client
}


func NewGithubClient(httpClient *http.Client) Github {
	client := github.NewClient(httpClient)

	return Github{Client: client}
}

// CreateTag adds a tag
func (gh *Github) CreateTag(name string, sha string) (*github.Reference, error) {
	ref := &github.Reference{
		Ref: github.String("refs/tags/" + name),
		Object: &github.GitObject{
			SHA: github.String(sha),
		},
	}

	ref, _, err := gh.Client.Git.CreateRef(context.Background(), owner, repo, ref)
	if err != nil {
		return nil, err
	}
	return ref, nil
}
