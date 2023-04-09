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
	case testIC1.frontMatter.URL:
		return &IssueContent{
			frontMatter: &YAMLFrontMatter{
				URL:    testIC1.frontMatter.URL,
				Title:  testIC1.frontMatter.Title,
				Labels: testIC1.frontMatter.Labels,
			},
			content: &Content{
				contents: testIC1.content.contents,
			},
		}, nil
	case testIC2.frontMatter.URL:
		return &IssueContent{
			frontMatter: &YAMLFrontMatter{
				URL:    testIC2.frontMatter.URL,
				Title:  testIC2.frontMatter.Title,
				Labels: testIC2.frontMatter.Labels,
			},
			content: &Content{
				contents: testIC2.content.contents,
			},
		}, nil
	}
	return nil, errors.New("unexpected url")
}

type mockFailGitHubClient struct{}

func NewMockFailGitHubClient(config config.Config) GitHubClient {
	return &mockFailGitHubClient{}
}

func (m *mockFailGitHubClient) GetIssueContent(issueURL string) (*IssueContent, error) {
	switch issueURL {
	case testIC1.frontMatter.URL:
		return &IssueContent{
			frontMatter: &YAMLFrontMatter{
				URL:    testIC1.frontMatter.URL,
				Title:  "different title for failure",
				Labels: testIC1.frontMatter.Labels,
			},
			content: &Content{
				contents: testIC1.content.contents,
			},
		}, nil
	}
	return nil, errors.New("unexpected url")
}
