package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client provides git operations
type Client interface {
	// HasDiff checks if there are any changes in the given directory
	HasDiff(dir string) (bool, error)
	// Commit commits all changes
	Commit(dir, message string) error
	// CommitAndPush commits all changes and pushes to remote
	CommitAndPush(dir, message string) error
}

type gitClient struct {
	userEmail string
	userName  string
}

// NewGitClient creates a new git client
func NewGitClient(userEmail, userName string) Client {
	return &gitClient{
		userEmail: userEmail,
		userName:  userName,
	}
}

// HasDiff checks if there are any changes in the working directory
func (gc *gitClient) HasDiff(dir string) (bool, error) {
	// Add untracked files to the index without staging them (git add -N .)
	cmd := exec.Command("git", "add", "-N", ".")
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return false, fmt.Errorf("git add -N: %w, output: %s", err, string(output))
	}

	// Check if there are any differences (git diff --name-only --exit-code)
	// Exit code 0: no diff, Exit code 1: has diff
	cmd = exec.Command("git", "diff", "--name-only", "--exit-code")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Exit code 1 means there are differences
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return true, nil
		}
		return false, fmt.Errorf("git diff: %w, output: %s", err, string(output))
	}
	// Exit code 0 means no differences
	return false, nil
}

// Commit commits all changes
func (gc *gitClient) Commit(dir, message string) error {
	// Configure git user
	if err := gc.configUser(dir); err != nil {
		return fmt.Errorf("config user: %w", err)
	}

	// Add all changes
	cmd := exec.Command("git", "add", "--all")
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add --all: %w, output: %s", err, string(output))
	}

	// Commit changes
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit: %w, output: %s", err, string(output))
	}

	return nil
}

// CommitAndPush commits all changes and pushes to remote
func (gc *gitClient) CommitAndPush(dir, message string) error {
	// Commit first
	if err := gc.Commit(dir, message); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	// Push to remote
	cmd := exec.Command("git", "push")
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push: %w, output: %s", err, string(output))
	}

	return nil
}

// configUser configures git user email and name
func (gc *gitClient) configUser(dir string) error {
	// Set user.email
	cmd := exec.Command("git", "config", "--local", "user.email", gc.userEmail)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config user.email: %w, output: %s", err, string(output))
	}

	// Set user.name
	cmd = exec.Command("git", "config", "--local", "user.name", gc.userName)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config user.name: %w, output: %s", err, string(output))
	}

	return nil
}

// ShowDiff shows git diff output (for debugging)
func (gc *gitClient) ShowDiff(dir string) (string, error) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
