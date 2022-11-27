---
title: Multi stage build binary
url: https://github.com/go-zen-chu/issue2md/issues/9
labels: [enhancement,good first issue]
---
Why
- Current docker image requires building go binary before building docker image

Proposal
- multi stage docker image and build go binary inside dockerfile so you don't need to build it during CI
solved by https://github.com/go-zen-chu/issue2md/commit/359d19b293e108342756cdb2f85a0b3da2ccb827
