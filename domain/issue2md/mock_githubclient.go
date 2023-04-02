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
	case TestIC1.frontMatter.url:
		return &IssueContent{
			frontMatter: &YAMLFrontMatter{
				url:    TestIC1.frontMatter.url,
				title:  TestIC1.frontMatter.title,
				labels: TestIC1.frontMatter.labels,
			},
			content: &Content{
				contents: TestIC1.content.contents,
			},
		}, nil
	}
	return nil, errors.New("unexpected url")
}
