package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-zen-chu/issue2md/internal/config"
	ui2m "github.com/go-zen-chu/issue2md/usecase/issue2md"
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
		genGitHubClient func(c config.Config) ui2m.GitHubClient
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
				genGitHubClient: ui2m.NewMockGitHubClient,
			},
			false,
		},
		{
			"When cmd args given, it should be successful when a markdown with same issue url but same title exists",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"test_token",
					"-export-dir",
					expDirAbsPath,
				},
				genGitHubClient: ui2m.NewMockGitHubClient,
			},
			false,
		},
		{
			"When same issue url with different title is trying to be exported, it should fail",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"test_token",
					"-export-dir",
					expDirAbsPath,
				},
				genGitHubClient: ui2m.NewMockFailGitHubClient,
			},
			true,
		},
		{
			`When both env vars and cmd args given, it should be successful and cmd args overwrite env vars values.
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
				genGitHubClient: ui2m.NewMockGitHubClient,
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
				genGitHubClient: ui2m.NewMockGitHubClient,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := setup(tt.args.envVars, tt.args.cmdArgs)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("setup() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if cfg == nil {
				// help message was shown
				if tt.wantErr {
					t.Errorf("setup() returned nil config, wantErr %v", tt.wantErr)
				}
				return
			}
			ghClient := tt.args.genGitHubClient(cfg)
			if err := run(cfg, ghClient); (err != nil) != tt.wantErr {
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
		genGitHubClient func(c config.Config) ui2m.GitHubClient
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
				genGitHubClient: ui2m.NewMockGitHubClient,
			},
			false,
		},
		{
			"When valid relative path given in cmd args, it should be successful",
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
				genGitHubClient: ui2m.NewMockGitHubClient,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := setup(tt.args.envVars, tt.args.cmdArgs)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("setup() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if cfg == nil {
				// help message was shown
				if tt.wantErr {
					t.Errorf("setup() returned nil config, wantErr %v", tt.wantErr)
				}
				return
			}
			ghClient := tt.args.genGitHubClient(cfg)
			if err := run(cfg, ghClient); (err != nil) != tt.wantErr {
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
		genGitHubClient func(c config.Config) ui2m.GitHubClient
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
				genGitHubClient: ui2m.NewMockGitHubClient,
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
				genGitHubClient: ui2m.NewMockGitHubClient,
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
				genGitHubClient: ui2m.NewMockGitHubClient,
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
				genGitHubClient: ui2m.NewMockGitHubClient,
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
				genGitHubClient: ui2m.NewMockGitHubClient,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := setup(tt.args.envVars, tt.args.cmdArgs)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("setup() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if cfg == nil {
				// help message was shown
				if tt.wantErr {
					t.Errorf("setup() returned nil config, wantErr %v", tt.wantErr)
				}
				return
			}
			ghClient := tt.args.genGitHubClient(cfg)
			if err := run(cfg, ghClient); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
