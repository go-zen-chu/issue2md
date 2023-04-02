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
	expDirAbsPath := filepath.Join(tmpDir, "issue")
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
					"ISSUE2MD_GITHUB_TOKEN=invalid_token_for_test",
					"ISSUE2MD_EXPORT_DIR=" + expDirAbsPath,
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			false,
		},
		{
			"When invalid env vars given, it should fail",
			args{
				envVars: []string{
					"ISSUE2MD_GITHUB_ISSUE_URL=https://github.com/Codertocat/Hello-World/issues/1",
					"ISSUE2MD_GITHUB_TOKEN=invalid_token_for_test",
					"ISSUE2MD_THIS_IS_INVALID=hehehe",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
		{
			"When both env vars and cmd args given, it should work and cmd args overwrite env vars values",
			args{
				envVars: []string{
					"ISSUE2MD_GITHUB_ISSUE_URL=overwritten",
					"ISSUE2MD_GITHUB_TOKEN=overwritten",
					"ISSUE2MD_THIS_IS_INVALID=overwritten",
				},
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"invalid_token_for_test",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
		{
			"When directory traversed path given, it should reject",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"invalid_token_for_test",
					"-export-dir",
					"./../root/",
				},
				genGitHubClient: di2m.NewMockGitHubClient,
			},
			true,
		},
		{
			"When valid relative path given, it should work",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/Codertocat/Hello-World/issues/1",
					"-github-token",
					"invalid_token_for_test",
					"-export-dir",
					"./",
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
					"invalid_token_for_test",
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
