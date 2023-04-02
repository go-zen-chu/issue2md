package issue2md

import (
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

// Load only YAML front matter for memory efficiency
func LoadFrontMatterFromMarkdownFile(filePath string) (*YAMLFrontMatter, error) {
	bt, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	matches := markdownRegex.FindSubmatch(bt)
	if matches != nil {
		return nil, fmt.Errorf("could not find front matter in file: %s", filePath)
	}
	var yfm *YAMLFrontMatter
	if err := yaml.Unmarshal(matches[1], yfm); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %s\n%s", filePath, matches[1])
	}
	return yfm, nil
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

func LoadFromMarkdownFile(filePath string) (*IssueContent, error) {
	bt, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	matches := markdownRegex.FindSubmatch(bt)
	if matches != nil {
		return nil, fmt.Errorf("could not find front matter in file: %s", filePath)
	}
	var yfm *YAMLFrontMatter
	if err := yaml.Unmarshal(matches[1], yfm); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %s\n%s", filePath, matches[1])
	}
	return &IssueContent{
		frontMatter: yfm,
		content: &Content{
			contents: strings.Split(string(matches[2]), "\n"),
		},
	}, nil
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
