package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	envPrefix         = "ISSUE2MD_"
	envExportDir      = envPrefix + "EXPORT_DIR"
	envGitHubIssueURL = envPrefix + "GITHUB_ISSUE_URL"
	envGitHubToken    = envPrefix + "GITHUB_TOKEN"
)

var (
	// singleton config
	cnf           *config
	edVal         *string
	ghIssueUrlVal *string
	ghTokenVal    *string
	checkDupsVal  *bool
)

// Config create interface for handling application config
type Config interface {
	LoadFromEnvVars(envVars []string) error
	SetupCommandArgs(flgSet *flag.FlagSet)
	LoadFromCommandArgs(flagName string) error
	GetExportDir() string
	GetGitHubIssueURL() string
	GetGitHubToken() string
	GetCheckDups() bool
}

type config struct {
	exportDirPath *exportDirPath
	ghIssueUrl    string
	ghToken       string
	checkDups     bool
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
		switch evs[0] {
		case envExportDir:
			edp, err := newExportDirPath(evs[1])
			if err != nil {
				return err
			}
			c.exportDirPath = edp
		case envGitHubIssueURL:
			c.ghIssueUrl = evs[1]
		case envGitHubToken:
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
	edVal = flgSet.String("export-dir", "./", fmt.Sprintf("Target directory to export issue as markdowns. Default is repository root ./ (%s)", envExportDir))
	ghIssueUrlVal = flgSet.String("issue-url", "", fmt.Sprintf("Set GitHub issue url (%s)", envGitHubIssueURL))
	ghTokenVal = flgSet.String("github-token", "", fmt.Sprintf("[WARN: recommended set from envvar %s] Set GitHub Token (%s)", envGitHubToken, envGitHubToken))
	checkDupsVal = flgSet.Bool("check-dups", false, "Find duplicate issueURL markdowns in export-dir")
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
	case "check-dups":
		c.checkDups = *checkDupsVal
	}
	return nil
}

func (c *config) GetExportDir() string {
	return c.exportDirPath.absPath
}

func (c *config) GetGitHubIssueURL() string {
	return c.ghIssueUrl
}

func (c *config) GetGitHubToken() string {
	return c.ghToken
}

func (c *config) GetCheckDups() bool {
	return c.checkDups
}
