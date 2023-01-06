---
title: Support **not** to create markdown when issue closed with `staled` or `dupes`
url: https://github.com/go-zen-chu/issue2md/issues/7
labels: [enhancement,good first issue]
---
Why
- github issues supports close reason: https://github.blog/changelog/2022-05-19-the-new-github-issues-may-19th-update/#%F0%9F%95%B5%F0%9F%8F%BD%E2%99%80%EF%B8%8F-issue-closed-reasons
- When you accidentally create two issues with same content, you want to close one of them with `duplicated`

Proposal
- By default, issue2md should **not** create markdown file when issue except `closed(completed)`
You can see sample payload when issue event fired

https://docs.github.com/ja/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#webhook-payload-example-when-someone-edits-an-issue


```json
[
  {
    "id": 1,
    "node_id": "MDEwOklzc3VlRXZlbnQx",
    "url": "https://api.github.com/repos/octocat/Hello-World/issues/events/1",
    "actor": {
      "login": "octocat",
      .
    },
    "event": "closed",
    "commit_id": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
    "commit_url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
    "created_at": "2011-04-14T16:00:49Z",
    "issue": {
      "id": 1,
      ……
      "number": 1347,
      "state": "open",
      "title": "Found a bug",
      "body": "I'm having a problem with this.",
      "user": {
        "login": "octocat",
        ……
      },
      "labels": [
        {
          "id": 208045946,
          "node_id": "MDU6TGFiZWwyMDgwNDU5NDY=",
          "url": "https://api.github.com/repos/octocat/Hello-World/labels/bug",
          "name": "bug",
          "description": "Something isn't working",
          "color": "f29513",
          "default": true
        }
      ],
      "assignee": {
        
      },
      "assignees": [
        {
          "login": "octocat",
          ……
        }
      ],
      "milestone": {
        ……
        "creator": {

        },
        "open_issues": 4,
        "closed_issues": 8,
        "created_at": "2011-04-10T20:09:31Z",
        "updated_at": "2014-03-03T18:58:10Z",
        "closed_at": "2013-02-12T13:22:01Z",
        "due_on": "2012-10-09T23:39:01Z"
      },
      "locked": true,
      "active_lock_reason": "too heated",
      "comments": 0,
      "pull_request": {
        "url": "https://api.github.com/repos/octocat/Hello-World/pulls/1347",
        "html_url": "https://github.com/octocat/Hello-World/pull/1347",
        "diff_url": "https://github.com/octocat/Hello-World/pull/1347.diff",
        "patch_url": "https://github.com/octocat/Hello-World/pull/1347.patch"
      },
      ……
      "state_reason": "completed"   ←----  when closed, you got state_reason
    }
  }
]
```
