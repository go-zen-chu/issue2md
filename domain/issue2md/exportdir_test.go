package issue2md

import (
	"reflect"
	"testing"
)

func TestNewExportDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *ExportDir
		wantErr bool
	}{
		// TBD
		// {
		// 	"if valid arg given, it should work",
		// 	args{
		// 		"issues/",
		// 	},
		// 	&ExportDir{
		// 		argPath: "issues/",
		// 	},
		// 	false,
		// },
		{
			"if invalid path given, it should get error",
			args{
				"../../root/",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExportDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExportDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExportDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExportDir_GetAbsPath(t *testing.T) {
	type fields struct {
		argPath string
		absPath string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := &ExportDir{
				argPath: tt.fields.argPath,
				absPath: tt.fields.absPath,
			}
			if got := ed.GetAbsPath(); got != tt.want {
				t.Errorf("ExportDir.GetAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
