package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	envPrefix = "ISSUE2MD_"
)

var (
	cnf           *config
	edVal         *string
	ghIssueUrlVal *string
	ghTokenVal    *string
)

// Config create interface for handling application config
type Config interface {
	LoadFromEnvVars(envVars []string) error
	SetupCommandArgs(flgSet *flag.FlagSet)
	LoadFromCommandArgs(flagName string) error
	GetExportDir() string
	GetGitHubIssueURL() string
	GetGitHubToken() string
}

type config struct {
	exportDirPath *exportDirPath
	ghIssueUrl    string
	ghToken       string
}

// validate and check relative or absolute path
type exportDirPath struct {
	givenPath string
	absPath   string
}

func newExportDirPath(givenPath string) (*exportDirPath, error) {
	edp := &exportDirPath{}
	if strings.HasPrefix(givenPath, "..") || strings.Contains(givenPath, "/../") {
		return nil, fmt.Errorf("traversing dir is not allowed: %s", givenPath)
	}
	edp.givenPath = givenPath
	if !filepath.IsAbs(givenPath) {
		// check whether relative path
		var wd string
		var err error
		if wd, err = os.Getwd(); err != nil {
			return nil, err
		}
		absPath := filepath.Join(wd, givenPath)
		if _, err := os.Stat(absPath); err != nil {
			return nil, fmt.Errorf("no such path: %s", absPath)
		}
		edp.absPath = absPath
	} else {
		// if path exists, it's absolute path
		edp.absPath = givenPath
	}
	return edp, nil
}

// NewConfig return globally singleton config
func NewConfig() Config {
	if cnf == nil {
		cnf = &config{}
	}
	return cnf
}

func (c *config) String() string {
	ghToken := "<empty>"
	if len(c.ghToken) > 0 {
		// make sure NOT TO print credentials
		ghToken = "<masked>"
	}
	return fmt.Sprintf("config{exportDir:%s,issueURL:%s,token:%s}", c.exportDirPath, c.ghIssueUrl, ghToken)
}

func (c *config) LoadFromEnvVars(envVars []string) error {
	var err error
	for _, ev := range envVars {
		evs := strings.SplitN(ev, "=", 2)
		if !strings.HasPrefix(evs[0], envPrefix) {
			continue
		}
		evKey := strings.TrimPrefix(evs[0], envPrefix)
		switch evKey {
		case "EXPORT_DIR":
			edp, err := newExportDirPath(evs[1])
			if err != nil {
				return err
			}
			c.exportDirPath = edp
		case "GITHUB_ISSUE_URL":
			c.ghIssueUrl = evs[1]
		case "GITHUB_TOKEN":
			c.ghToken = evs[1]
		default:
			err = fmt.Errorf("%w: invalid env var %s", err, ev)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *config) SetupCommandArgs(flgSet *flag.FlagSet) {
	edVal = flgSet.String("export-dir", "./", "Target directory to export issue as markdowns. Default is repository root")
	ghIssueUrlVal = flgSet.String("issue-url", "", fmt.Sprintf("Set GitHub issue url (%sGITHUB_ISSUE_URL)", envPrefix))
	ghTokenVal = flgSet.String("github-token", "", fmt.Sprintf("[WARN: recommended set from envvar %sGITHUB_TOKEN] Set GitHub Token (%sGITHUB_TOKEN)", envPrefix, envPrefix))
}

func (c *config) LoadFromCommandArgs(flagName string) error {
	switch flagName {
	case "export-dir":
		edp, err := newExportDirPath(*edVal)
		if err != nil {
			return err
		}
		c.exportDirPath = edp
	case "issue-url":
		c.ghIssueUrl = *ghIssueUrlVal
	case "github-token":
		c.ghToken = *ghTokenVal
	}
	return nil
}

func (c *config) GetExportDir() string {
	return cnf.exportDirPath.absPath
}

func (c *config) GetGitHubIssueURL() string {
	return cnf.ghIssueUrl
}

func (c *config) GetGitHubToken() string {
	return cnf.ghToken
}
