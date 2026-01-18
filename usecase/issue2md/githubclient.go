package issue2md

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

type GitHubClient interface {
	GetIssueContent(issueURL string) (*IssueContent, error)
}
