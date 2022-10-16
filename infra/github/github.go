package github

import (
	dis "github.com/go-zen-chu/issue2md/domain/issue2md"
)

type ghClient struct {
}

func NewGitHubClient() dis.GitHubClient {
	return &ghClient{}
}

func (ghc *ghClient) GetIssueInfo() error {
	return nil
}
