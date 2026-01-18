# issue2md

Simple GitHub Action for archiving GitHub issue to Markdown and store in git.

## Workflow

| kind           | status                                                                                                                                                                                                                                                                                                                                                                                                                |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Code Check     | [![Go general check (go fmt, vet, build, test)](https://github.com/go-zen-chu/issue2md/actions/workflows/go-ci.yaml/badge.svg)](https://github.com/go-zen-chu/issue2md/actions/workflows/go-ci.yaml)                                                                                                                                                                                                                  |
| E2E Test       | [![E2E Test in PR](https://github.com/go-zen-chu/issue2md/actions/workflows/pr-e2e-test.yaml/badge.svg)](https://github.com/go-zen-chu/issue2md/actions/workflows/pr-e2e-test.yaml)                                                                                                                                                                                                                                   |
| Docker Publish | [![Actions Status](https://github.com/go-zen-chu/issue2md/workflows/main-push-image/badge.svg)](https://github.com/go-zen-chu/issue2md/actions/workflows/main-push-image.yaml) <br> [dockerhub: amasuda/issue2md](https://hub.docker.com/repository/docker/amasuda/issue2md)                                                                                                                                                     |
| Dog fooding    | [![issue2md](https://github.com/go-zen-chu/issue2md/actions/workflows/main-issue2md.yaml/badge.svg)](https://github.com/go-zen-chu/issue2md/actions/workflows/main-issue2md.yaml) <br> [![issue2md-check-dups](https://github.com/go-zen-chu/issue2md/actions/workflows/main-issue2md-check-dups.yaml/badge.svg?branch=main)](https://github.com/go-zen-chu/issue2md/actions/workflows/main-issue2md-check-dups.yaml) |

## Goal

The goal of this project is,

* Convert GitHub issues to markdown files so you can manage & archive issues as markdown with git

### Use case

* Archiving issues to markdowns for permenently store them in git
  * Example: Write down your private notes to GitHub issues and save them as markdown files

## Installation

1. Run a command below in your repository root.

    ```bash
    mkdir -p .github/workflows; 
    curl -s https://raw.githubusercontent.com/go-zen-chu/issue2md/main/docs/issue2md.yml -o .github/workflows/issue2md.yml
    ```

1. Commit and push to your repository
1. When issue closed, GitHub action called `issue2md` export the issue as a markdown to specified dir. Please see an example [issue](https://github.com/go-zen-chu/issue2md/issues/2) & [archived issue](https://github.com/go-zen-chu/issue2md/blob/main/issues/test%20issue.md).

### Update issue2md action

Basically, run a command below to get latest action file.

```bash
curl -s https://raw.githubusercontent.com/go-zen-chu/issue2md/main/docs/issue2md.yml -o .github/workflows/issue2md.yml
```

If you edited upper file locally, then hit `git diff` and merge to the latest file.

## Parameters

| name         | value type | required | default | env var                   | description                                                                                                                                                                                                             |
| ------------ | ---------- | -------- | ------- | ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| github_token | string     | *        | -       | ISSUE2MD_GITHUB_TOKEN              | [GitHub token](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#using-the-github_token-in-a-workflow) for r/w issue and repo. Using `${{ secrets.GITHUB_TOKEN }}` as inputs satisfies. |
| issue-url    | string     | *        | -       | ISSUE2MD_GITHUB_ISSUE_URL | Issue url for exporting to markdown. Using `${{ github.event.issue.html_url }}` as inputs satisfies.                                                                                                                    |
| export-dir   | string     |          | `.`     | ISSUE2MD_EXPORT_DIR       | A directory in your repository to export markdown. Default is `.` (repository root)                                                                                                                                     |
| check-dups   | bool       |          | false   | ISSUE2MD_CHECK_DUPS       | Optional flag for checking duplicate issue URL markdown exists in export-dir. If there is any, issue2md print which files are duplicated.                                                                               |

> [!NOTE]
> Why `ISSUE2MD_GITHUB_TOKEN` and not `GITHUB_TOKEN`?
> `GITHUB_TOKEN` is used in command like `gh` which may affect when you want to separately manage tokens in the same action.

## Development

### test localy

You can test this action locally as follows.

```bash
source .envrc  # set ISSUE2MD_GITHUB_TOKEN env var
go run cmd/issue2md -debug -export-dir=./issues -issue-url=https://github.com/go-zen-chu/issue2md/issues/2

# or using nerdctl. Replace `lima nerdctl` to `docker` if you use docker
# build image
GOARCH=amd64 GOOS=linux go build -v -o issue2md ./cmd/issue2md; lima nerdctl build -t issue2md:latest .
# run on container
lima nerdctl run -it -e ISSUE2MD_GITHUB_TOKEN=${ISSUE2MD_GITHUB_TOKEN} --rm issue2md:latest -- -debug -issue-url=https://github.com/go-zen-chu/issue2md/issues/2 
```

## Appendix

### Markdown metadata format

Metadata of an issue are stored with [YAML Front-matter](https://jekyllrb.com/docs/front-matter/) format.
