package main

import (
	"os"
	"path/filepath"
	"testing"

	di2m "github.com/go-zen-chu/issue2md/domain/issue2md"
	"github.com/go-zen-chu/issue2md/internal/config"
)

func Test_run(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "issue2md-*")
	if err != nil {
		t.Errorf("mkdir temp: %s", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("remove all %s: %s", tmpDir, err)
		}
	}()
	expDirAbsPath := filepath.Join(tmpDir, "issues")
	if err := os.Mkdir(expDirAbsPath, 0777); err != nil {
		t.Errorf("mkdir %s: %s", expDirAbsPath, err)
	}
	// generated
	type args struct {
		envVars         []string
		cmdArgs         []string
		genGitHubClient func(c config.Config) di2m.GitHubClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"When only env vars given, it should work",
			args{
				envVars: []string{
					"ISSUE2MD_GITHUB_ISSUE_URL=https://github.com/Codertocat/Hello-World/issues/1",
					"ISSUE2MD_GITHUB_TOKEN=test_token",
					"ISSUE2MD_EXPORT_DIR=" + expDirAbsPath,
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			false,
		},
		{
			`When both env vars and cmd args given, it should work and cmd args overwrite env vars values.
				Export will succeed when exporting an issue that is not exists`,
			args{
				envVars: []string{
					"ISSUE2MD_GITHUB_ISSUE_URL=overwritten",
					"ISSUE2MD_GITHUB_TOKEN=overwritten",
					"ISSUE2MD_EXPORT_DIR=./",
				},
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/2",
					"-github-token",
					"test_token",
					"-export-dir",
					expDirAbsPath,
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			false,
		},
		{
			"When help flag given, show help and exit without error",
			args{
				cmdArgs: []string{
					"issue2md",
					"-help",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.envVars, tt.args.cmdArgs, tt.args.genGitHubClient); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test cases for relative path. We need to remove file exported in relative path after test performed.
func Test_run_relative_path(t *testing.T) {
	// remove md file generated to relative path
	defer func() {
		files, err := filepath.Glob("./*.md")
		if err != nil {
			t.Errorf("error while globbing files: %s", err)
		}
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				t.Errorf("remove %s: %s", file, err)
			}
		}
	}()
	type args struct {
		envVars         []string
		cmdArgs         []string
		genGitHubClient func(c config.Config) di2m.GitHubClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"When both env vars and cmd args given, it should work and cmd args overwrite env vars values",
			args{
				envVars: []string{
					"ISSUE2MD_GITHUB_ISSUE_URL=overwritten",
					"ISSUE2MD_GITHUB_TOKEN=overwritten",
					"ISSUE2MD_EXPORT_DIR=./",
				},
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"test_token",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			false,
		},
		{
			// this test check relative path, make sure no file or valid md file exists
			"When valid relative path given, it should get error because same issue file already exists",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"test_token",
					"-export-dir",
					"./",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.envVars, tt.args.cmdArgs, tt.args.genGitHubClient); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_run_failure_case(t *testing.T) {
	// generated
	type args struct {
		envVars         []string
		cmdArgs         []string
		genGitHubClient func(c config.Config) di2m.GitHubClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"When invalid env vars given, it should fail",
			args{
				envVars: []string{
					"ISSUE2MD_GITHUB_ISSUE_URL=https://github.com/Codertocat/Hello-World/issues/1",
					"ISSUE2MD_GITHUB_TOKEN=test_token",
					"ISSUE2MD_THIS_IS_INVALID=hehehe",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
		{
			"When a directory traversed path given, it should reject",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"test_token",
					"-export-dir",
					"./../root/",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
		{
			"When invalid path given, it should fail",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"test_token",
					"-export-dir",
					"/no_such_root_path/in_your_computer",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
		{
			"When invalid command flag given, it should fail",
			args{
				cmdArgs: []string{
					"issue2md",
					"-no-such-flag",
					"https://github.com/Codertocat/Hello-World/issues/1",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
		{
			"When debug flag given, log with debug mode and fails because of insufficient args",
			args{
				cmdArgs: []string{
					"issue2md",
					"-debug",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.envVars, tt.args.cmdArgs, tt.args.genGitHubClient); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
