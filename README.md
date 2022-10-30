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

You can test this action locally by,

```bash
go run cmd/issue2md/main.go -export-dir=issues -debug -issue-url=https://github.com/go-zen-chu/issue2md/issues/2
```

### package structure

1. domain: implements core logics
2. infra: implements interface of domain logic
3. internal: application wide logics

## Appendix

### Metadata format

Converted metadata of a issue are stored as [YAML Front-matter](https://jekyllrb.com/docs/front-matter/).
