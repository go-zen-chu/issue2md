package issue2md

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewIssueContent(t *testing.T) {
	type args struct {
		url      string
		title    string
		labels   []string
		contents []string
	}
	tests := []struct {
		name string
		args args
		want *IssueContent
	}{
		{
			"if valid args given, it should create new IssueContent",
			args{
				url:      testIC1.frontMatter.URL,
				title:    testIC1.frontMatter.Title,
				labels:   testIC1.frontMatter.Labels,
				contents: testIC1.content.contents,
			},
			&testIC1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIssueContent(tt.args.url, tt.args.title, tt.args.labels, tt.args.contents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIssueContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadFrontMatterFromMarkdownFile_LineEndings(t *testing.T) {
	t.Parallel()

	const (
		title  = "Reading book"
		url    = "https://github.com/test-org/test-repo/issues/11"
		labels = "[book]"
	)

	buildMarkdown := func(eol string, leadingBlankLine bool, sepWhitespace bool) string {
		openSep := frontMatterSep
		closeSep := frontMatterSep
		if sepWhitespace {
			openSep = "  " + frontMatterSep + "  "
			closeSep = frontMatterSep + "   "
		}

		yamlLines := []string{
			"title: \"" + title + "\"",
			"url: " + url,
			"labels: " + labels,
		}
		bodyLines := []string{
			"",
			"## Expected output",
			"- Something concrete",
		}

		var sb strings.Builder
		if leadingBlankLine {
			sb.WriteString(eol)
		}
		sb.WriteString(openSep + eol)
		sb.WriteString(strings.Join(yamlLines, eol) + eol)
		sb.WriteString(closeSep + eol)
		sb.WriteString(strings.Join(bodyLines, eol) + eol)
		return sb.String()
	}

	lf := buildMarkdown("\n", false, false)
	leadingBlankLF := buildMarkdown("\n", true, false)
	sepWithSpacesLF := buildMarkdown("\n", false, true)

	crlf := buildMarkdown("\r\n", false, false)
	bomCRLF := string([]byte{0xEF, 0xBB, 0xBF}) + crlf

	tests := []struct {
		name    string
		content string
	}{
		{name: "LF", content: lf},
		{name: "LeadingBlankLine+LF", content: leadingBlankLF},
		{name: "CRLF", content: crlf},
		{name: "UTF8BOM+CRLF", content: bomCRLF},
		{name: "SepWhitespace+LF", content: sepWithSpacesLF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			p := filepath.Join(dir, "test.md")
			if err := os.WriteFile(p, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("write test file: %v", err)
			}
			yfm, err := LoadFrontMatterFromMarkdownFile(p)
			if err != nil {
				t.Fatalf("LoadFrontMatterFromMarkdownFile() error = %v", err)
			}
			if yfm.Title != title {
				t.Fatalf("Title = %q, want %q", yfm.Title, title)
			}
			if yfm.URL != url {
				t.Fatalf("URL = %q, want %q", yfm.URL, url)
			}
			if !reflect.DeepEqual(yfm.Labels, []string{"book"}) {
				t.Fatalf("Labels = %#v, want %#v", yfm.Labels, []string{"book"})
			}
		})
	}
}

func TestIssueContent_GetMDFilename(t *testing.T) {
	type fields struct {
		url      string
		title    string
		labels   []string
		contents []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"if valid title given, it should work",
			fields{
				title: "test issue",
			},
			"test issue.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &IssueContent{
				frontMatter: &YAMLFrontMatter{
					URL:    tt.fields.url,
					Title:  tt.fields.title,
					Labels: tt.fields.labels,
				},
				content: &Content{
					contents: tt.fields.contents,
				},
			}
			if got := ic.GetMDFilename(); got != tt.want {
				t.Errorf("IssueContent.GetMDFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueContent_Print(t *testing.T) {
	type fields struct {
		url      string
		title    string
		labels   []string
		contents []string
	}
	type args struct {
		contentseparator string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"If valid IssueContent given, it should work",
			fields{
				url:      testIC1.frontMatter.URL,
				title:    testIC1.frontMatter.Title,
				labels:   testIC1.frontMatter.Labels,
				contents: testIC1.content.contents,
			},
			args{
				"\n",
			},
			testIC1Output,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &IssueContent{
				frontMatter: &YAMLFrontMatter{
					URL:    tt.fields.url,
					Title:  tt.fields.title,
					Labels: tt.fields.labels,
				},
				content: &Content{
					contents: tt.fields.contents,
				},
			}
			var got string
			var err error
			if got, err = ic.GenerateContent(tt.args.contentseparator); got != tt.want {
				t.Errorf("IssueContent.Print() = %v, want %v", got, tt.want)
			}
			if (err == nil) == tt.wantErr {
				t.Errorf("Expected wantErr: %t, but err is %s", tt.wantErr, err)
			}
		})
	}
}
