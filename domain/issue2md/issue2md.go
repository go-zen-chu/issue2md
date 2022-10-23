package issue2md

type Issue2md interface {
	Convert2md() error
}

type issue2md struct {
	ghClient GitHubClient
	ghi      *IssueContent
}

func NewIssue2md(ghClient GitHubClient, ghi *IssueContent) Issue2md {
	return &issue2md{
		ghClient: ghClient,
		ghi:      ghi,
	}
}

func (i2m *issue2md) Convert2md() error {
	// Do some process using ghi
	return nil
}
