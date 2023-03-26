package issue2md

import (
	"fmt"
	"os"
	"path/filepath"
)

type Issue2md interface {
	Convert2md(issueURL string) error
}

type issue2md struct {
	ghClient GitHubClient
	expDir   string
	ghi      *IssueContent
}

func NewIssue2md(ghClient GitHubClient, expDir string) Issue2md {
	return &issue2md{
		ghClient: ghClient,
		expDir:   expDir,
	}
}

func (i2m *issue2md) Convert2md(issueURL string) error {
	ic, err := i2m.ghClient.GetIssueContent(issueURL)
	if err != nil {
		return fmt.Errorf("get issue content: %w", err)
	}
	mdStr := ic.GenerateContent("\n")
	if err := os.WriteFile(filepath.Join(i2m.expDir, ic.GetMDFilename()), []byte(mdStr), 0755); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
