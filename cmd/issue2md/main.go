package main

import (
	"flag"
	"fmt"
	"os"

	dis "github.com/go-zen-chu/issue2md/domain/issue2md"
	igh "github.com/go-zen-chu/issue2md/infra/github"
	"github.com/go-zen-chu/issue2md/internal/log"
)

var (
	mainFlagSet   = flag.NewFlagSet("issue2md", flag.ExitOnError)
	debugVal      = mainFlagSet.Bool("debug", false, "Enable debug")
	helpVal       = mainFlagSet.Bool("help", false, "Show help")
	edVal         *exportDirValue
	ghIssueUrlVal = mainFlagSet.String("issue-url", "", "Set GitHub issue url")
)

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

func init() {
	edVal = new(exportDirValue)
	mainFlagSet.Var(edVal, "export-dir", "Target directory to export issue as markdowns. Default is '/' which is repository root")
}

func help() {
	fmt.Println("usage: issue2md <flags>")
	mainFlagSet.PrintDefaults()
}

func main() {
	if len(os.Args) == 1 {
		help()
		os.Exit(1)
	}
	if err := mainFlagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "parse args: %s", err)
		os.Exit(1)
	}
	if *helpVal {
		help()
		os.Exit(0)
	}
	if err := log.Init(*debugVal); err != nil {
		fmt.Fprintf(os.Stderr, "initializing logger: %s", err)
		os.Exit(1)
	}
	defer log.Close()
	log.Debugf("helpFlag:%t,debugFlag:%t,exportDir:%+v", *helpVal, *debugVal, *edVal.ed)

	ghClient := igh.NewGitHubClient()
	ghi := dis.NewIssueContent(*ghIssueUrlVal)
	log.Debugf("GitHub Issue: %+v", ghi)
	i2m := dis.NewIssue2md(ghClient, ghi)
	if err := i2m.Convert2md(); err != nil {
		log.Fatal(fmt.Sprintf("converting to markdown: %s", err))
	}
}
