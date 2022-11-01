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
	var wd string
	var err error
	if wd, err = os.Getwd(); err != nil {
		return nil, err
	}
	absPath := filepath.Join(wd, path)
	if _, err = os.Stat(absPath); err != nil {
		return nil, err
	}
	return &ExportDir{
		argPath: path,
		absPath: absPath,
	}, nil
}

func (ed *ExportDir) GetAbsPath() string {
	return ed.absPath
}
