package issue2md

// GitHubClient is domain repository interface for r/w data with GitHub
type GitHubClient interface {
	GetIssueContent(issueURL string) (*IssueContent, error)
}
