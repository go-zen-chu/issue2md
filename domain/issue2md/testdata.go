package issue2md

// somehow, TestIC1 defined in issuecontent_test.go cannot be accessed from other file in same package.
var (
	TestIC1 = IssueContent{
		frontMatter: &YAMLFrontMatter{
			url:    "https://github.com/Codertocat/Hello-World/issues/1",
			title:  "test issue",
			labels: []string{"a", "b"},
		},
		content: &Content{
			contents: []string{"test1", "test2"},
		},
	}
)

const (
	TestIC1Output = `---
title: test issue
url: https://github.com/Codertocat/Hello-World/issues/1
labels: [a,b]
---
test1
test2
`
)
