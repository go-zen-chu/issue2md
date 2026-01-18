package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config create interface for handling application config
type Config interface {
	LoadFromEnvVars(envVars []string) error
	LoadFromCommandArgs(args []string) error
	GetExportDir() string
	GetGitHubIssueURL() string
	GetGitHubToken() string
	GetCheckDups() bool
	IsDebugMode() bool
	ShowHelp() (bool, string)
}

const (
	envPrefix         = "ISSUE2MD_"
	envExportDir      = envPrefix + "EXPORT_DIR"
	envGitHubIssueURL = envPrefix + "GITHUB_ISSUE_URL"
	envGitHubToken    = "GITHUB_TOKEN"
)

type config struct {
	exportDirPath *exportDirPath
	ghIssueUrl    string
	ghToken       string
	checkDups     bool
	debug         bool
	help          bool
	flgSet        *flag.FlagSet
}

// NewConfig return globally singleton config
func NewConfig() Config {
	return &config{
		flgSet: flag.NewFlagSet("issue2md", flag.ContinueOnError),
	}
}

// validate and check relative or absolute path
type exportDirPath struct {
	givenPath string
	absPath   string
}

func (c *config) LoadFromCommandArgs(args []string) error {
	if len(args) <= 1 {
		// no args given
		return nil
	}
	if c.flgSet.Parsed() {
		// already parsed
		return nil
	}
	edVal := c.flgSet.String("export-dir", "./", fmt.Sprintf("Target directory to export issue as markdowns. Default is repository root ./ (%s)", envExportDir))
	ghIssueUrlVal := c.flgSet.String("issue-url", "", fmt.Sprintf("Set GitHub issue url (%s)", envGitHubIssueURL))
	ghTokenVal := c.flgSet.String("github-token", "", fmt.Sprintf("[WARN: recommended set from envvar %s] Set GitHub Token (%s)", envGitHubToken, envGitHubToken))
	checkDupsVal := c.flgSet.Bool("check-dups", false, "Find duplicate issueURL markdowns in export-dir")
	debugVal := c.flgSet.Bool("debug", false, "Enable debug")
	helpVal := c.flgSet.Bool("help", false, "Show help")

	var errg error
	if err := c.flgSet.Parse(args[1:]); err != nil {
		return fmt.Errorf("parse args: %s", err)
	}
	c.flgSet.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "export-dir":
			edp, err := newExportDirPath(*edVal)
			if err != nil {
				errg = errors.Join(errg, fmt.Errorf("handling flag %s: %w", f.Name, err))
			} else {
				c.exportDirPath = edp
			}
		case "issue-url":
			c.ghIssueUrl = *ghIssueUrlVal
		case "github-token":
			c.ghToken = *ghTokenVal
		case "check-dups":
			c.checkDups = *checkDupsVal
		case "debug":
			c.debug = *debugVal
		case "help":
			c.help = *helpVal
		default:
			errg = errors.Join(errg, fmt.Errorf("handling flag %s: unknown flag", f.Name))
		}
	})
	// Apply default export-dir even when flag wasn't explicitly provided.
	if c.exportDirPath == nil {
		edp, err := newExportDirPath(*edVal)
		if err != nil {
			errg = errors.Join(errg, fmt.Errorf("handling flag export-dir: %w", err))
		} else {
			c.exportDirPath = edp
		}
	}
	return errg
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

func (c *config) LoadFromEnvVars(envVars []string) error {
	for _, ev := range envVars {
		evs := strings.SplitN(ev, "=", 2)
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
		}
	}
	// Apply default export-dir when env var wasn't provided.
	if c.exportDirPath == nil {
		edp, err := newExportDirPath("./")
		if err != nil {
			return err
		}
		c.exportDirPath = edp
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

func (c *config) IsDebugMode() bool {
	return c.debug
}

func (c *config) ShowHelp() (bool, string) {
	var sb strings.Builder
	sb.WriteString("Usage of issue2md:\n")
	c.flgSet.VisitAll(func(f *flag.Flag) {
		sb.WriteString(fmt.Sprintf("  -%s: %s\n", f.Name, f.Usage))
	})
	return c.help, sb.String()
}

func (c *config) String() string {
	ghToken := "<empty>"
	if len(c.ghToken) > 0 {
		// make sure NOT TO print credentials
		ghToken = "<masked>"
	}
	return fmt.Sprintf("config{exportDir:%s,issueURL:%s,token:%s}", c.exportDirPath, c.ghIssueUrl, ghToken)
}
