# issue2md

[![Actions Status](https://github.com/go-zen-chu/issue2md/workflows/ci/badge.svg)](https://github.com/go-zen-chu/issue2md/actions/workflows/ci.yml)
[![Actions Status](https://github.com/go-zen-chu/issue2md/workflows/push-image/badge.svg)](https://github.com/go-zen-chu/issue2md/actions/workflows/push-image.yml)
[![Actions Status](https://github.com/go-zen-chu/issue2md/workflows/issue2md/badge.svg)](https://github.com/go-zen-chu/issue2md/actions/workflows/issue2md.yml)

[dockerhub: amasuda/issue2md](https://hub.docker.com/repository/docker/amasuda/issue2md)

GitHub Action for archiving GitHub issue to Markdown

## Goal

The goal of this project is,

- Convert GitHub issues to markdowns so that you can manage & archive by git

### Use case

- Archiving all issues to markdowns so that you can permenently store them in git
- Write down your private note to issues and save them as markdown

## How to use

### GitHub Action yaml

TBD

### Parameters

TBD

## Development

### test localy

You can test this action locally by, setting ISSUE2MD_GITHUB_TOKEN envvar and run

```bash
# With go
go run cmd/issue2md/*.go -debug -export-dir=./issues -issue-url=https://github.com/go-zen-chu/issue2md/issues/2

# or using nerdctl. Replace `lima nerdctl` to `docker` if you use docker
# build image
GOARCH=amd64 GOOS=linux go build -v -o issue2md ./cmd/issue2md/*.go; lima nerdctl build -t issue2md:latest .
# run on container
lima nerdctl run -it -e ISSUE2MD_GITHUB_TOKEN=${ISSUE2MD_GITHUB_TOKEN} --rm issue2md:latest -- -debug -issue-url=https://github.com/go-zen-chu/issue2md/issues/2 
```

### package structure

1. domain: implements core logics
2. infra: implements interface of domain logic
3. internal: application wide logics

## Appendix

### Metadata format

Converted metadata of a issue are stored as [YAML Front-matter](https://jekyllrb.com/docs/front-matter/).
