package main

import (
	"fmt"
	"os"

	"github.com/go-zen-chu/issue2md/infra/github"
	"github.com/go-zen-chu/issue2md/internal/config"
	"github.com/go-zen-chu/issue2md/internal/log"
	ui2m "github.com/go-zen-chu/issue2md/usecase/issue2md"
)

func main() {
	if err := run(os.Environ(), os.Args); err != nil {
		panic(fmt.Sprintf("run(): %s", err))
	}
}

func run(envVars []string, cmdArgs []string) error {
	cfg := config.NewConfig()
	if err := cfg.LoadFromEnvVars(envVars); err != nil {
		return fmt.Errorf("load env var: %w", err)
	}
	if err := cfg.LoadFromCommandArgs(cmdArgs); err != nil {
		return fmt.Errorf("load args: %w", err)
	}
	if help, helpMsg := cfg.ShowHelp(); help {
		fmt.Println(helpMsg)
		return nil
	}
	// initialize logger
	if err := log.Init(cfg.IsDebugMode()); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer log.Close()
	log.Debugf("config: %+v", cfg)

	githubClient := github.NewGitHubClient(github.Token(cfg.GetGitHubToken()))

	i2muc := ui2m.NewIssue2mdUseCase(githubClient, cfg.GetExportDir())
	if cfg.GetCheckDups() {
		res, err := i2muc.CheckDuplicateIssueFile()
		if len(res) == 0 {
			res = "No duplicate issueURL markdown files :tada:"
		}
		fmt.Printf("[CheckDuplicateIssueURLFile] %s\n", res)
		if err != nil {
			fmt.Printf("error: %s", err)
		}
	} else {
		if err := i2muc.Convert2md(c.GetGitHubIssueURL()); err != nil {
			return fmt.Errorf("convert to markdown: %w", err)
		}
		log.Infof("Export issue %s to %s, succeeded\n", cfg.GetGitHubIssueURL(), cfg.GetExportDir())
	}
	return nil
}
