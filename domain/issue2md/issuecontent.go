package issue2md

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	controlCharsRegex     = regexp.MustCompile("[\u0000-\u001f\u0080-\u009f]")
	relativePathRegex     = regexp.MustCompile(`^\.+`)
	filenameReservedRegex = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
)

type IssueContent struct {
	url      string
	title    string
	labels   []string
	comments []string
}

func NewIssueContent(url, title string, labels []string, comments []string) *IssueContent {
	return &IssueContent{
		url:      url,
		title:    title,
		labels:   labels,
		comments: comments,
	}
}

func (ic *IssueContent) GetMDFilename() string {
	return ic.title + ".md"
}

func (ic *IssueContent) GenerateContent(commentSeparator string) string {
	var sb strings.Builder
	// YAML front matter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", ic.title))
	sb.WriteString(fmt.Sprintf("url: %s\n", ic.url))
	sb.WriteString(fmt.Sprintf("labels: [%s]\n", strings.Join(ic.labels, ",")))
	sb.WriteString("---\n")
	// content
	sb.WriteString(strings.Join(ic.comments, commentSeparator))
	sb.WriteString("\n")
	return sb.String()
}
