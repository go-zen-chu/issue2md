package issue2md

import (
	"fmt"
	"strings"
)

type IssueContent struct {
	url      string
	title    string
	labels   []string
	contents []string
}

func NewIssueContent(url, title string, labels []string, contents []string) *IssueContent {
	return &IssueContent{
		url:      url,
		title:    title,
		labels:   labels,
		contents: contents,
	}
}

func (ic *IssueContent) GetMDFilename() string {
	return ic.title + ".md"
}

func (ic *IssueContent) GenerateContent(contentseparator string) string {
	var sb strings.Builder
	// YAML front matter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", ic.title))
	sb.WriteString(fmt.Sprintf("url: %s\n", ic.url))
	sb.WriteString(fmt.Sprintf("labels: [%s]\n", strings.Join(ic.labels, ",")))
	sb.WriteString("---\n")
	// content
	sb.WriteString(strings.Join(ic.contents, contentseparator))
	sb.WriteString("\n")
	return sb.String()
}
