package issue2md

// NOTICE: somehow, testIC1 defined in issuecontent_test.go cannot be accessed from other file in same package.
var (
	testIC1 = IssueContent{
		frontMatter: &YAMLFrontMatter{
			URL:    "https://github.com/Codertocat/Hello-World/issues/1",
			Title:  "test issue",
			Labels: []string{"a", "b"},
		},
		content: &Content{
			contents: []string{"test1", "test2"},
		},
	}
)

const (
	testIC1Output = `---
title: test issue
url: https://github.com/Codertocat/Hello-World/issues/1
labels: [a, b]
---
test1
test2
`
)
