// NOTICE: this package must be imported from test files, not production codes.
package testdata

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetRepositoryRootPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current dir: %w", err)
	}
	repoRoot := cwd
	for {
		if _, err := os.Stat(filepath.Join(repoRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(repoRoot)
		if parent == repoRoot {
			return "", fmt.Errorf("go.mod not found in ancestors of %s", cwd)
		}
		repoRoot = parent
	}
	return repoRoot, nil
}
