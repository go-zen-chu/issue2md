package issue2md

import (
	"reflect"
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
