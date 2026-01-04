package issue2md

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	frontMatterSep = "---"
)

var (
	utf8BOMString = "\ufeff"
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

func extractYAMLFrontMatterFromReader(r io.Reader) ([]byte, bool, error) {
	scanner := bufio.NewScanner(r)
	// Default token limit (64K) can be too small for unusual YAML lines.
	// We only scan until the closing '---', but be defensive.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	foundStart := false
	var yamlLines []string
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if !foundStart {
			line = strings.TrimPrefix(line, utf8BOMString)
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			if trimmed != frontMatterSep {
				// No front matter in the beginning
				return nil, false, nil
			}
			foundStart = true
			continue
		}

		if strings.TrimSpace(line) == frontMatterSep {
			return []byte(strings.Join(yamlLines, "\n")), true, nil
		}
		yamlLines = append(yamlLines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, false, fmt.Errorf("scan file for front matter extraction: %w", err)
	}
	return nil, false, nil
}

// Load only YAML front matter for memory efficiency
func LoadFrontMatterFromMarkdownFile(filePath string) (*YAMLFrontMatter, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = f.Close() }()

	yamlBytes, ok, err := extractYAMLFrontMatterFromReader(f)
	if err != nil {
		return nil, fmt.Errorf("extract front matter: %w", err)
	}
	if !ok {
		if len(yamlBytes) == 0 {
			return nil, fmt.Errorf("could not find front matter in file: %s", filePath)
		}
		return nil, fmt.Errorf("could not find closing front matter separator in file: %s\nfront matter:\n%s", filePath, string(yamlBytes))
	}
	var yfm YAMLFrontMatter
	if err := yaml.Unmarshal(yamlBytes, &yfm); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %s: %w\n%s", filePath, err, string(yamlBytes))
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
