package main

import (
	"fmt"
	"os"

	di2m "github.com/go-zen-chu/issue2md/domain/issue2md"
	"github.com/go-zen-chu/issue2md/internal/config"
	"github.com/go-zen-chu/issue2md/internal/di"
	"github.com/go-zen-chu/issue2md/internal/log"
	"github.com/go-zen-chu/issue2md/internal/runner"
	ui2m "github.com/go-zen-chu/issue2md/usecase/issue2md"
)

func run(
	envVars []string,
	cmdArgs []string,
	genGitHubClient func(c config.Config) di2m.GitHubClient) error {
	r := runner.NewRunner("issue2md")
	if err := r.LoadConfigFromEnvVars(envVars); err != nil {
		return fmt.Errorf("load env var: %w", err)
	}
	if err := r.LoadConfigFromCommandArgs(cmdArgs); err != nil {
		return fmt.Errorf("load args: %w", err)
	}
	if err := r.Run(func(c config.Config) error {
		ghClient := genGitHubClient(c)
		i2muc := ui2m.NewIssue2mdUseCase(ghClient, c.GetExportDir())
		if err := i2muc.Convert2md(c.GetGitHubIssueURL()); err != nil {
			return fmt.Errorf("convert to markdown: %w", err)
		}
		log.Infof("Export issue %s to %s, succeeded\n", c.GetGitHubIssueURL(), c.GetExportDir())
		return nil
	}); err != nil {
		return fmt.Errorf("while running: %w", err)
	}
	return nil
}

func main() {
	di := di.NewDIContainer()
	if err := run(os.Environ(), os.Args, di.GitHubClient); err != nil {
		panic(fmt.Sprintf("run(): %s", err))
	}
}
