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
	mainFlagSet = flag.NewFlagSet("issue2md", flag.ExitOnError)
	debugFlag   = mainFlagSet.Bool("debug", false, "Enable debug")
	helpFlag    = mainFlagSet.Bool("help", false, "Show help")
	edValue     *exportDirValue
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
	edValue = &exportDirValue{}
	mainFlagSet.Var(edValue, "export-dir", "Directory for exporting markdown. Default is '/' which is repository root")
}

func help() {
	fmt.Println("usage: issue2md <flags>")
	mainFlagSet.PrintDefaults()
}

func main() {
	if len(os.Args) == 1 {
		help()
		os.Exit(0)
	}
	if err := mainFlagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "parse args: %s", err)
		os.Exit(1)
	}
	if *helpFlag {
		help()
		os.Exit(0)
	}
	if err := log.Init(*debugFlag); err != nil {
		fmt.Fprintf(os.Stderr, "initializing logger: %s", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Infof("helpFlag:%t,debugFlag:%t,exportDir:%+v", *helpFlag, *debugFlag, *edValue.ed)

	ghClient := igh.NewGitHubClient()
	i2m := dis.NewIssue2md(ghClient)
	if err := i2m.Convert2md(); err != nil {
		log.Fatal(fmt.Sprintf("converting to markdown: %s", err))
	}
}
