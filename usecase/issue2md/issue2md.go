package issue2md

import (
	di2m "github.com/go-zen-chu/issue2md/domain/issue2md"
)

type Issue2mdUseCase interface {
	Convert2md(issueURL string) error
}

type issue2mdUsecase struct {
	i2m di2m.Issue2md
}

func NewIssue2mdUseCase(ghClient di2m.GitHubClient, expDir string) Issue2mdUseCase {
	return &issue2mdUsecase{
		i2m: di2m.NewIssue2md(ghClient, expDir),
	}
}

// Usecase Convert2md convert issue to markdown
func (imu *issue2mdUsecase) Convert2md(issueURL string) error {
	return imu.i2m.Convert2md(issueURL)
}

// If more usecases required, should be implemented below
