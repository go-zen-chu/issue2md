package issue2md

import (
	"errors"

	"github.com/go-zen-chu/issue2md/internal/config"
)

type mockGitHubClient struct{}

func NewMockGitHubClient(config config.Config) GitHubClient {
	return &mockGitHubClient{}
}

func (m *mockGitHubClient) GetIssueContent(issueURL string) (*IssueContent, error) {
	switch issueURL {
	case TestIC1.frontMatter.URL:
		return &IssueContent{
			frontMatter: &YAMLFrontMatter{
				URL:    TestIC1.frontMatter.URL,
				Title:  TestIC1.frontMatter.Title,
				Labels: TestIC1.frontMatter.Labels,
			},
			content: &Content{
				contents: TestIC1.content.contents,
			},
		}, nil
	}
	return nil, errors.New("unexpected url")
}
