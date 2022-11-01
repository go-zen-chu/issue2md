package issue2md

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-zen-chu/issue2md/internal/log"
)

type Issue2md interface {
	Convert2md(issueURL string) error
}

type issue2md struct {
	ghClient GitHubClient
	expDir   *ExportDir
	ghi      *IssueContent
}

func NewIssue2md(ghClient GitHubClient, expDir *ExportDir) Issue2md {
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
	log.Debugf("generated content\n%s", mdStr)
	if err := os.WriteFile(filepath.Join(i2m.expDir.absPath, ic.GetMDFilename()), []byte(mdStr), 0755); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
