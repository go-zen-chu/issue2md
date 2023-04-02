//go:build wireinject
// +build wireinject

package di

import (
	di2m "github.com/go-zen-chu/issue2md/domain/issue2md"
	igh "github.com/go-zen-chu/issue2md/infra/github"
	"github.com/go-zen-chu/issue2md/internal/config"
	"github.com/go-zen-chu/issue2md/internal/runner"
	"github.com/google/wire"
)

func initRunner(appName string) runner.Runner {
	panic(wire.Build(runner.NewRunner))
}

func toOptionSetterSlice(os igh.OptionSetter) []igh.OptionSetter {
	return []igh.OptionSetter{os}
}

func getConfigGitHubIssueToken(c config.Config) string {
	return c.GetGitHubToken()
}

func initGitHubClient(c config.Config) di2m.GitHubClient {
	panic(wire.Build(igh.NewGitHubClient, toOptionSetterSlice, igh.Token, getConfigGitHubIssueToken))
}
