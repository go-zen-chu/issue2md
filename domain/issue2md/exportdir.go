package issue2md

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	repoRootRegexp = regexp.MustCompile(`.+?issue2md`)
)

type ExportDir struct {
	argPath string
	absPath string
}

func NewExportDir(path string) (*ExportDir, error) {
	if strings.HasPrefix(path, "..") {
		return nil, fmt.Errorf("traversing dir is not allowed: %s", path)
	}
	r, err := getRepoRoot()
	if err != nil {
		return nil, err
	}
	return &ExportDir{
		argPath: path,
		absPath: filepath.Join(r, path),
	}, nil
}

func getRepoRoot() (string, error) {
	f, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working dir: %w", err)
	}
	repoRoot := repoRootRegexp.FindString(f)
	if len(repoRoot) == 0 {
		return "", fmt.Errorf("could not find repo root in path: %s", f)
	}
	return repoRoot, nil
}

func (ed *ExportDir) GetAbsPath() string {
	return ed.absPath
}
