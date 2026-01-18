package main

import (
	"fmt"
	"os"

	"github.com/go-zen-chu/issue2md/infra/git"
	"github.com/go-zen-chu/issue2md/infra/github"
	"github.com/go-zen-chu/issue2md/internal/config"
	"github.com/go-zen-chu/issue2md/internal/log"
	ui2m "github.com/go-zen-chu/issue2md/usecase/issue2md"
)

func main() {
	var cfg config.Config
	var err error
	if cfg, err = setup(os.Environ(), os.Args); err != nil {
		panic(fmt.Sprintf("setup: %s", err))
	}
	if cfg == nil {
		// help message was shown
		return
	}
	if err := run(cfg, github.NewGitHubClient(github.Token(cfg.GetGitHubToken()))); err != nil {
		panic(fmt.Sprintf("run: %s", err))
	}
}

func setup(envVars []string, cmdArgs []string) (config.Config, error) {
	cfg := config.NewConfig()
	if err := cfg.LoadFromEnvVars(envVars); err != nil {
		return nil, fmt.Errorf("load env var: %w", err)
	}
	if err := cfg.LoadFromCommandArgs(cmdArgs); err != nil {
		return nil, fmt.Errorf("load args: %w", err)
	}
	if help, helpMsg := cfg.ShowHelp(); help {
		fmt.Println(helpMsg)
		return nil, nil
	}
	// initialize logger
	if err := log.Init(cfg.IsDebugMode()); err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	defer log.Close()
	log.Debugf("config: %+v", cfg)

	return cfg, nil
}

func run(cfg config.Config, githubClient ui2m.GitHubClient) error {
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
		if err := i2muc.Convert2md(cfg.GetGitHubIssueURL()); err != nil {
			return fmt.Errorf("convert to markdown: %w", err)
		}
		log.Infof("Export issue %s to %s, succeeded\n", cfg.GetGitHubIssueURL(), cfg.GetExportDir())

		// Auto commit and push if enabled
		if cfg.IsAutoCommit() || cfg.IsAutoPush() {
			if err := handleGitOperations(cfg); err != nil {
				return fmt.Errorf("git operations: %w", err)
			}
		}
	}
	return nil
}

func handleGitOperations(cfg config.Config) error {
	gitClient := git.NewGitClient("action@github.com", "GitHub Action")
	
	// Check if there are any changes
	hasDiff, err := gitClient.HasDiff(cfg.GetExportDir())
	if err != nil {
		return fmt.Errorf("check diff: %w", err)
	}

	if !hasDiff {
		log.Infof("No changes detected, skipping git operations")
		return nil
	}

	log.Infof("Changes detected in export directory")

	// Commit changes if auto-commit is enabled
	if cfg.IsAutoCommit() {
		commitMsg := "[skip ci] [GitHub Action] Update automatically"
		
		if cfg.IsAutoPush() {
			// Commit and push
			log.Infof("Committing and pushing changes...")
			if err := gitClient.CommitAndPush(cfg.GetExportDir(), commitMsg); err != nil {
				return fmt.Errorf("commit and push: %w", err)
			}
			log.Infof("Changes committed and pushed successfully")
		} else {
			// Commit only (no push)
			log.Infof("Committing changes...")
			if err := gitClient.Commit(cfg.GetExportDir(), commitMsg); err != nil {
				return fmt.Errorf("commit: %w", err)
			}
			log.Infof("Changes committed successfully")
		}
	}

	return nil
}
