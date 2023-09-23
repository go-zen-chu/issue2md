---
title: Dependency Dashboard
url: https://github.com/go-zen-chu/issue2md/issues/24
labels: []
---
This issue lists Renovate updates and detected dependencies. Read the [Dependency Dashboard](https://docs.renovatebot.com/key-concepts/dashboard/) docs to learn more.

## Open

These updates have all been created already. Click a checkbox below to force a retry/rebase of any.

 - [ ] <!-- rebase-branch=renovate/actions-checkout-4.x -->[Update actions/checkout action to v4](../pull/27)
 - [ ] <!-- rebase-branch=renovate/actions-setup-go-4.x -->[Update actions/setup-go action to v4](../pull/28)
 - [ ] <!-- rebase-branch=renovate/mr-smithers-excellent-docker-build-push-6.x -->[Update mr-smithers-excellent/docker-build-push action to v6](../pull/29)
 - [ ] <!-- rebase-all-open-prs -->**Click on this checkbox to rebase all open PRs at once**

## Detected dependencies

<details><summary>dockerfile</summary>
<blockquote>

<details><summary>Dockerfile</summary>

 - `golang 1.20.5-buster`

</details>

</blockquote>
</details>

<details><summary>github-actions</summary>
<blockquote>

<details><summary>.github/workflows/ci.yml</summary>

 - `actions/checkout v3`
 - `actions/setup-go v3`

</details>

<details><summary>.github/workflows/issue2md-check-dups.yml</summary>

 - `actions/checkout v3`
 - `slackapi/slack-github-action v1.24.0`

</details>

<details><summary>.github/workflows/issue2md.yml</summary>

 - `actions/checkout v3`

</details>

<details><summary>.github/workflows/push-image.yml</summary>

 - `actions/checkout v3`
 - `mr-smithers-excellent/docker-build-push v5`
 - `mr-smithers-excellent/docker-build-push v5`

</details>

<details><summary>.github/workflows/test-issue2md.yml</summary>

 - `actions/checkout v3`
 - `slackapi/slack-github-action v1.24.0`

</details>

</blockquote>
</details>

<details><summary>gomod</summary>
<blockquote>

<details><summary>go.mod</summary>

 - `go 1.19`
 - `github.com/google/go-github v17.0.0+incompatible`
 - `github.com/google/wire v0.5.0`
 - `go.uber.org/zap v1.26.0`
 - `golang.org/x/oauth2 v0.12.0`
 - `gopkg.in/yaml.v3 v3.0.1`

</details>

</blockquote>
</details>

---

- [ ] <!-- manual job -->Check this box to trigger a request for Renovate to run again on this repository


