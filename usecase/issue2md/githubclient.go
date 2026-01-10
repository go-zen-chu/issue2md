package issue2md

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

type GitHubClient interface {
	GetIssueContent(issueURL string) (*IssueContent, error)
}
