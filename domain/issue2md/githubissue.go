package issue2md

type GitHubIssue struct {
	number int
	title  string
	label  string
}

func NewGitHubIssue(number int, title string, label string) *GitHubIssue {
	return &GitHubIssue{
		number: number,
		title:  title,
		label:  label,
	}
}
