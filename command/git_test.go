package command

import "testing"

func TestGitCommand_Basic(t *testing.T) {
	if gitCommand.Use != "git" {
				// git command should have use "git"
		t.Fatalf("expected git command use 'git', got %s", gitCommand.Use)
	}
}
