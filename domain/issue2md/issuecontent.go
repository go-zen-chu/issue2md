package issue2md

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	frontMatterSep = "---"
)

var (
	markdownRegex = regexp.MustCompile(`(?ms)^---\n(.+?)\n---\n(.*)$`)
)

type IssueContent struct {
	frontMatter *YAMLFrontMatter
	content     *Content
}

type YAMLFrontMatter struct {
	Title  string   `yaml:"title"`
	URL    string   `yaml:"url"`
	Labels []string `yaml:"labels,flow"`
}

type Content struct {
	contents []string
}

func firstNLines(bt []byte, n int) string {
	if n <= 0 || len(bt) == 0 {
		return ""
	}
	parts := bytes.SplitN(bt, []byte("\n"), n+1)
	if len(parts) > n {
		parts = parts[:n]
	}
	return string(bytes.Join(parts, []byte("\n")))
}

// Load only YAML front matter for memory efficiency
func LoadFrontMatterFromMarkdownFile(filePath string) (*YAMLFrontMatter, error) {
	bt, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	matches := markdownRegex.FindSubmatch(bt)
	if matches == nil {
		return nil, fmt.Errorf("could not find front matter in file: %s\nfirst 10 lines:\n%s", filePath, firstNLines(bt, 10))
	}
	var yfm YAMLFrontMatter
	if err := yaml.Unmarshal(matches[1], &yfm); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %s\n%s", filePath, matches[1])
	}
	return &yfm, nil
}

func (yfm *YAMLFrontMatter) GetIssueURL() string {
	return yfm.URL
}

func NewIssueContent(url, title string, labels []string, contents []string) *IssueContent {
	return &IssueContent{
		frontMatter: &YAMLFrontMatter{
			URL:    url,
			Title:  title,
			Labels: labels,
		},
		content: &Content{
			contents: contents,
		},
	}
}

func (ic *IssueContent) GetMDFilename() string {
	return ic.frontMatter.Title + ".md"
}

func (ic *IssueContent) GenerateContent(contentseparator string) (string, error) {
	fm, err := yaml.Marshal(ic.frontMatter)
	if err != nil {
		return "", fmt.Errorf("marshal yaml: %w", err)
	}
	var sb strings.Builder
	sb.WriteString(frontMatterSep + "\n")
	sb.WriteString(string(fm))
	sb.WriteString(frontMatterSep + "\n")
	sb.WriteString(strings.Join(ic.content.contents, contentseparator))
	sb.WriteString("\n")
	return sb.String(), nil
}
