package issue2md

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	di2m "github.com/go-zen-chu/issue2md/domain/issue2md"
)

type Issue2mdUseCase interface {
	Convert2md(issueURL string) error
	CheckDuplicateIssueFile() (string, error)
}

type issue2mdUsecase struct {
	i2m    di2m.Issue2md
	expDir string
}

func NewIssue2mdUseCase(ghClient di2m.GitHubClient, expDir string) Issue2mdUseCase {
	return &issue2mdUsecase{
		i2m:    di2m.NewIssue2md(ghClient, expDir),
		expDir: expDir,
	}
}

// Usecase for converting github issue to markdown
func (imu *issue2mdUsecase) Convert2md(issueURL string) error {
	return imu.i2m.Convert2md(issueURL)
}

type fileMeta struct {
	finfo       fs.FileInfo
	frontMatter *di2m.YAMLFrontMatter
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
		yfm, err := di2m.LoadFrontMatterFromMarkdownFile(absPath)
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
