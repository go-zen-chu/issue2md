package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	dis "github.com/go-zen-chu/issue2md/domain/issue2md"
)

const (
	appName   = "issue2md"
	envPrefix = "ISSUE2MD_"
)

var (
	flgSet        = flag.NewFlagSet(appName, flag.ExitOnError)
	debugVal      = flgSet.Bool("debug", false, "Enable debug")
	helpVal       = flgSet.Bool("help", false, "Show help")
	edVal         *exportDirValue
	ghIssueUrlVal = flgSet.String("issue-url", "", fmt.Sprintf("Set GitHub issue url (%sGITHUB_ISSUE_URL)", envPrefix))
	ghTokenVal    = flgSet.String("github-token", "", fmt.Sprintf("[WARN: recommended set from envvar %sGITHUB_TOKEN] Set GitHub Token (%sGITHUB_TOKEN)", envPrefix))
)

func HelpString() string {
	var sb strings.Builder
	sb.WriteString("usage: ")
	sb.WriteString(appName)
	sb.WriteString(" <flags>\n")
	op := flgSet.Output()
	flgSet.SetOutput(&sb)
	// print to string builder
	flgSet.PrintDefaults()
	flgSet.SetOutput(op)
	return sb.String()
}

type config struct {
	debug      bool
	help       bool
	exportDir  *dis.ExportDir
	ghIssueUrl string
	ghToken    string
}

func NewConfig() *config {
	edVal = new(exportDirValue)
	flgSet.Var(edVal, "export-dir", "Target directory to export issue as markdowns. Default is '/' which is repository root")
	// return default config
	return &config{
		debug:      false,
		help:       false,
		exportDir:  nil,
		ghIssueUrl: "",
		ghToken:    "",
	}
}

func (c *config) LoadEnvVars(envVars []string) error {
	var err error
	for _, ev := range envVars {
		evs := strings.SplitN(ev, "=", 2)
		if !strings.HasPrefix(evs[0], envPrefix) {
			continue
		}
		evKey := strings.TrimPrefix(evs[0], envPrefix)
		switch evKey {
		case "GITHUB_ISSUE_URL":
			c.ghIssueUrl = evs[1]
		case "GITHUB_TOKEN":
			c.ghToken = evs[1]
		default:
			err = fmt.Errorf("%w: invalid env var %s", err, ev)
		}
	}
	return err
}

func (c *config) LoadCommandArgs(args []string) error {
	if !flgSet.Parsed() {
		if len(args) <= 1 {
			return fmt.Errorf("invalid args len: %d", len(os.Args))
		}
		if err := flgSet.Parse(args[1:]); err != nil {
			return fmt.Errorf("parse args: %s", err)
		}
		// visit specified flag
		flgSet.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "debug":
				c.debug = *debugVal
			case "help":
				c.help = *helpVal
			case "export-dir":
				c.exportDir = edVal.ed
			case "issue-url":
				c.ghIssueUrl = *ghIssueUrlVal
			case "github-token":
				c.ghToken = *ghTokenVal
			}
		})
	}
	return nil
}

type exportDirValue struct {
	ed *dis.ExportDir
}

// implements Value interface for flag argument
func (edv *exportDirValue) String() string {
	if edv == nil || edv.ed == nil {
		return "/"
	}
	return edv.ed.GetAbsPath()
}

// implements Value interface for flag argument
func (edv *exportDirValue) Set(path string) error {
	ed, err := dis.NewExportDir(path)
	if err != nil {
		return fmt.Errorf("invalid arg: %w", err)
	}
	edv.ed = ed
	return nil
}
