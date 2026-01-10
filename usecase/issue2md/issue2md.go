package issue2md

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	di2m "github.com/go-zen-chu/issue2md/domain/issue2md"
	"github.com/go-zen-chu/issue2md/infra/github"
)

type Issue2mdUseCase interface {
	Convert2md(issueURL string) error
	CheckDuplicateIssueFile() (string, error)
}

type issue2mdUsecase struct {
	ghClient github.GitHubClient
	expDir string
}

func NewIssue2mdUseCase(ghClient github.GitHubClient, expDir string) Issue2mdUseCase {
	return &issue2mdUsecase{
		ghClient: ghClient,
		expDir:   expDir,
	}
}

// Usecase for converting github issue to markdown
func (imu *issue2mdUsecase) Convert2md(issueURL string) error {
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

type fileMeta struct {
	finfo       fs.FileInfo
	frontMatter *YAMLFrontMatter
	absPath     string
}

// Usecase for finding duplicate file in export-dir
func (imu *issue2mdUsecase) CheckDuplicateIssueFile() (string, error) {
	files, err := os.ReadDir(imu.expDir)
	if err != nil {
		return "", fmt.Errorf("read dirs: %w", err)
	}
	var errg error
	var sb strings.Builder
	fileDict := make(map[string]*fileMeta, len(files))
	for _, file := range files {
		fi, err := file.Info()
		if err != nil {
			errg = fmt.Errorf("%w\nfile info:%s", errg, err)
			continue
		}
		if fi.IsDir() {
			continue
		}
		if !strings.HasSuffix(fi.Name(), ".md") {
			continue
		}
		absPath := filepath.Join(imu.expDir, fi.Name())
		yfm, err := LoadFrontMatterFromMarkdownFile(absPath)
		if err != nil {
			errg = fmt.Errorf("%w\nparse markdown:%s", errg, err)
			continue
		}
		iurl := yfm.GetIssueURL()
		if _, ok := fileDict[iurl]; !ok {
			fileDict[iurl] = &fileMeta{
				finfo:       fi,
				frontMatter: yfm,
				absPath:     absPath,
			}
			// unique issue markdown
			continue
		}
		sb.WriteString("Find duplicate issue files:\n")
		sb.WriteString(fileDict[iurl].absPath)
		sb.WriteString(", [Modified] ")
		sb.WriteString(fileDict[iurl].finfo.ModTime().String())
		sb.WriteString("\n")
		sb.WriteString(absPath)
		sb.WriteString(", [Modified] ")
		sb.WriteString(fi.ModTime().String())
		sb.WriteString("\n")
	}
	return sb.String(), errg
}
