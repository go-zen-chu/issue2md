package main

import (
	"fmt"
	"os"

	dis "github.com/go-zen-chu/issue2md/domain/issue2md"
	igh "github.com/go-zen-chu/issue2md/infra/github"
	"github.com/go-zen-chu/issue2md/internal/log"
)

func main() {
	conf := NewConfig()
	if err := conf.LoadEnvVars(os.Environ()); err != nil {
		panic(fmt.Errorf("load env var: %w", err))
	}
	if err := conf.LoadCommandArgs(os.Args); err != nil {
		panic(fmt.Errorf("load args: %w", err))
	}
	if conf.help {
		fmt.Println(HelpString())
		os.Exit(0)
	}
	if err := log.Init(conf.debug); err != nil {
		panic(fmt.Errorf("initializing logger: %w", err))
	}
	defer log.Close()
	log.Debugf("config: %+v", conf)

	//TODO: use di factory for future work
	ghClient := igh.NewGitHubClient(igh.Token(conf.ghToken))
	i2m := dis.NewIssue2md(ghClient, conf.exportDir)
	if err := i2m.Convert2md(conf.ghIssueUrl); err != nil {
		log.Fatal(fmt.Sprintf("converting to markdown: %s", err))
	}
	log.Infof("Export issue %s to %s, succeeded\n", conf.ghIssueUrl, conf.exportDir.GetAbsPath())
}
