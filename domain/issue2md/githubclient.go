package issue2md

type GitHubClient interface {
	GetIssueContent(issueURL string) (*IssueContent, error)
}
