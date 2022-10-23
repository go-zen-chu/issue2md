package github

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	dis "github.com/go-zen-chu/issue2md/domain/issue2md"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ghClient struct {
	ctx    context.Context
	option *Option
	cli    *github.Client
}

func NewGitHubClient(setters ...OptionSetter) dis.GitHubClient {
	o := &Option{
		baseURL: "https://github.com",
		token:   "",
	}
	for _, setter := range setters {
		setter(o)
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: o.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli := github.NewClient(tc)
	return &ghClient{
		ctx:    ctx,
		option: o,
		cli:    cli,
	}
}

func (ghc *ghClient) GetIssueContent(issueURL string) (*dis.IssueContent, error) {
	u, err := url.Parse(issueURL)
	if err != nil {
		return nil, fmt.Errorf("parse issueURL: %w", err)
	}
	if !strings.HasPrefix(issueURL, path.Join()) {
		return nil, fmt.Errorf("invalid url, baseURL: %s, given URL: %s", ghc.option.baseURL, issueURL)
	}
	// e.g. /Codertocat/Hello-World/issues/12
	p := u.Path
	ps := strings.Split(p, "/")
	if len(ps) != 5 {
		return nil, fmt.Errorf("invalid url. given URL: %s", issueURL)
	}
	owner := ps[1]
	repo := ps[2]
	number, err := strconv.Atoi(ps[4])
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", issueURL)
	}
	i, _, err := ghc.cli.Issues.Get(ghc.ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("get issue: %w", err)
	}
	labels := make([]string, len(i.Labels))
	for idx, l := range i.Labels {
		labels[idx] = *l.Name
	}
	cs, _, err := ghc.cli.Issues.ListComments(ghc.ctx, owner, repo, number, nil)
	if err != nil {
		return nil, fmt.Errorf("get comments: %w", err)
	}
	comments := make([]string, len(cs))
	for idx, c := range cs {
		comments[idx] = *c.Body
	}
	ic := dis.NewIssueContent(issueURL, i.GetTitle(), labels, comments)
	return ic, nil
}
