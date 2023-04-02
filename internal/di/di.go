package di

import (
	di2m "github.com/go-zen-chu/issue2md/domain/issue2md"
	"github.com/go-zen-chu/issue2md/internal/config"
	"github.com/go-zen-chu/issue2md/internal/runner"
)

// DI interface is defined for mocking DI container while testing
type DI interface {
	Runner(appName string) runner.Runner
	GitHubClient(c config.Config) di2m.GitHubClient
}

type di struct {
}

func NewDIContainer() DI {
	return &di{}
}

func (d *di) Runner(appName string) runner.Runner {
	return initRunner(appName)
}

func (d *di) GitHubClient(c config.Config) di2m.GitHubClient {
	return initGitHubClient(c)
}

func (d *di) MockGitHubClient(c config.Config) di2m.GitHubClient {
	return di2m.NewMockGitHubClient(c)
}
