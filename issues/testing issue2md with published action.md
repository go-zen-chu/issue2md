---
title: testing issue2md with published action
url: https://github.com/go-zen-chu/issue2md/issues/6
labels: []
---
This is a test issue for testing published action
## Notes

Without test label in this issue, test-issue2md action skips it process.

https://github.com/go-zen-chu/issue2md/actions/runs/3378734003

![image](https://user-images.githubusercontent.com/1454332/199526928-ac6547bd-b61b-42ae-9fb8-e27cfedcfe5d.png)

if test label given, test-issue2md action should work
as tested below, issue2md was skipped when `test` label given
![image](https://user-images.githubusercontent.com/1454332/199656264-add2a1af-af59-4e1d-ae59-2bf2adbc0024.png)

git diff returns empty string when a new file just created.
if there is no labels, issue2md action should run
