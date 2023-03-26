package main

import "testing"

func Test_run(t *testing.T) {
	type args struct {
		envVars []string
		cmdArgs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"When env vars given, it should fail due to test github token",
			args{
				envVars: []string{
					"ISSUE2MD_GITHUB_ISSUE_URL=https://github.com/go-zen-chu/issue2md/issues/15",
					"ISSUE2MD_GITHUB_TOKEN=invalid_token_for_test",
					"ISSUE2MD_EXPORT_DIR=issues",
				},
			},
			true,
		},
		{
			"When invalid env vars given, it should fail",
			args{
				envVars: []string{
					"ISSUE2MD_GITHUB_ISSUE_URL=https://github.com/go-zen-chu/issue2md/issues/15",
					"ISSUE2MD_GITHUB_TOKEN=invalid_token_for_test",
					"ISSUE2MD_THIS_IS_INVALID=hehehe",
				},
			},
			true,
		},
		{
			"When cmd args given, it should fail due to test github token",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/go-zen-chu/issue2md/issues/15",
					"-github-token",
					"invalid_token_for_test",
				},
			},
			true,
		},
		{
			"When directory traversed path given, it should reject",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/go-zen-chu/issue2md/issues/15",
					"-github-token",
					"invalid_token_for_test",
					"-export-dir",
					"../../root/",
				},
			},
			true,
		},
		{
			"When valid relative path given, it should work but fail due to test github token",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/go-zen-chu/issue2md/issues/15",
					"-github-token",
					"invalid_token_for_test",
					"-export-dir",
					"./",
				},
			},
			true,
		},
		{
			"When invalid root path given, it should fail",
			args{
				cmdArgs: []string{
					"issue2md",
					"-issue-url",
					"https://github.com/go-zen-chu/issue2md/issues/15",
					"-github-token",
					"invalid_token_for_test",
					"-export-dir",
					"/no_such_root_path/in_your_computer",
				},
			},
			true,
		},
		{
			"When flag given, it should fail",
			args{
				cmdArgs: []string{
					"issue2md",
					"-no-such-flag",
					"https://github.com/go-zen-chu/issue2md/issues/15",
				},
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
			},
			true,
		},
		{
			"When no arg given, it should fail",
			args{
				cmdArgs: []string{
					"issue2md",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.envVars, tt.args.cmdArgs); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
