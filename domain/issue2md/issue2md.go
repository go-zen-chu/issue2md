package issue2md

type Issue2md interface {
	Convert2md() error
}

type GitHubClient interface {
	GetIssueInfo() error
}

type issue2md struct {
	ghClient GitHubClient
}

func NewIssue2md(ghClient GitHubClient) Issue2md {
	return &issue2md{
		ghClient: ghClient,
	}
}

func (i2m *issue2md) Convert2md() error {
	return nil
}
