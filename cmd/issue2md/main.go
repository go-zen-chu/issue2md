package main

import (
	"fmt"
	"os"

	dis "github.com/go-zen-chu/issue2md/domain/issue2md"
	igh "github.com/go-zen-chu/issue2md/infra/github"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "initializing zap logger: %s", err)
		os.Exit(1)
	}
	defer logger.Sync()
	ghClient := igh.NewGitHubClient()
	i2m := dis.NewIssue2md(ghClient)
	if err = i2m.Convert2md(); err != nil {
		logger.Panic(fmt.Sprintf("converting to markdown: %s", err))
	}
}
