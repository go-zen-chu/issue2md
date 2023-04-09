package issue2md

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Issue2md interface {
	Convert2md(issueURL string) error
}

type issue2md struct {
	ghClient GitHubClient
	expDir   string
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
	// return error when duplicate file already exists
	files, err := os.ReadDir(i2m.expDir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", i2m.expDir, err)
	}
	for _, file := range files {
		fi, err := file.Info()
		if err != nil {
			return fmt.Errorf("checking file info (%s): %w", file.Name(), err)
		}
		if fi.IsDir() {
			continue
		}
		if !strings.HasSuffix(fi.Name(), ".md") {
			continue
		}
		absPath := filepath.Join(i2m.expDir, fi.Name())
		yfm, err := LoadFrontMatterFromMarkdownFile(absPath)
		if err != nil {
			return fmt.Errorf("could not load file %s: %w", absPath, err)
		}
		if yfm.GetIssueURL() == issueURL && yfm.Title != ic.frontMatter.Title {
			return fmt.Errorf("markdown with same issueURL (%s) but different title %s found: %s", issueURL, yfm.Title, absPath)
		}
	}
	mdStr, err := ic.GenerateContent("\n")
	if err != nil {
		return fmt.Errorf("generate content: %w", err)
	}
	if err := os.WriteFile(filepath.Join(i2m.expDir, ic.GetMDFilename()), []byte(mdStr), 0755); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
