package issue2md

type GitHubClient interface {
	GetIssueInfo() error
}
