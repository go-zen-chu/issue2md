// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/go-zen-chu/issue2md/domain/issue2md"
	"github.com/go-zen-chu/issue2md/infra/github"
	"github.com/go-zen-chu/issue2md/internal/config"
	"github.com/go-zen-chu/issue2md/internal/runner"
)

// Injectors from wire.go:

func initRunner(appName string) runner.Runner {
	runnerRunner := runner.NewRunner(appName)
	return runnerRunner
}

func initGitHubClient(c config.Config) issue2md.GitHubClient {
	string2 := getConfigGitHubIssueToken(c)
	optionSetter := github.Token(string2)
	v := toOptionSetterSlice(optionSetter)
	gitHubClient := github.NewGitHubClient(v...)
	return gitHubClient
}

// wire.go:

func toOptionSetterSlice(os github.OptionSetter) []github.OptionSetter {
	return []github.OptionSetter{os}
}

func getConfigGitHubIssueToken(c config.Config) string {
	return c.GetGitHubToken()
}
