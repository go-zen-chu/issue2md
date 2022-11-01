package issue2md

import (
	"reflect"
	"testing"
)

var (
	testIC1 = IssueContent{
		url:      "https://github.com/Codertocat/Hello-World/issues/1",
		title:    "test issue",
		labels:   []string{"a", "b"},
		contents: []string{"test1", "test2"},
	}
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
				url:      "https://github.com/Codertocat/Hello-World/issues/1",
				title:    "test issue",
				labels:   []string{"a", "b"},
				contents: []string{"test1", "test2"},
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
			"if valid fields given, it should work",
			fields{
				title: "test issue",
			},
			"test issue.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &IssueContent{
				url:      tt.fields.url,
				title:    tt.fields.title,
				labels:   tt.fields.labels,
				contents: tt.fields.contents,
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
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"if valid IssueContent given, it should work",
			fields{
				url:      "https://github.com/Codertocat/Hello-World/issues/1",
				title:    "test issue",
				labels:   []string{"a", "b"},
				contents: []string{"test1", "test2"},
			},
			args{
				"\n",
			},
			`---
title: test issue
url: https://github.com/Codertocat/Hello-World/issues/1
labels: [a,b]
---
test1
test2
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &IssueContent{
				url:      tt.fields.url,
				title:    tt.fields.title,
				labels:   tt.fields.labels,
				contents: tt.fields.contents,
			}
			if got := ic.GenerateContent(tt.args.contentseparator); got != tt.want {
				t.Errorf("IssueContent.Print() = %v, want %v", got, tt.want)
			}
		})
	}
}
