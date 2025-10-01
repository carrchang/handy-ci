package execution

import (
	"testing"

	"github.com/carrchang/handy-ci/config"
)

func TestGitExecution_CheckArgs_AlwaysNil(t *testing.T) {
	g := GitExecution{}
	if err := g.CheckArgs(nil, []string{"anything"}); err != nil {
		// GitExecution.CheckArgs currently never errors
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestGitExecution_Parse_Normal(t *testing.T) {
	workspace := config.Workspace{Name: "ws", Path: "/tmp/ws"}
	group := config.Group{Name: "group"}
	repo := config.Repository{Name: "repo"}

	cmd := fakeCobraCommand("git")

	execs, err := GitExecution{}.Parse(cmd, []string{"status"}, workspace, group, repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(execs) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(execs))
	}
	ex := execs[0]
	if ex.Command != "git" || len(ex.Args) != 1 || ex.Args[0] != "status" {
		t.Fatalf("unexpected execution: %+v", ex)
	}
	expectedPath := RepositoryPath(workspace, group, repo)
	if ex.Path != expectedPath {
		t.Fatalf("expected path %s, got %s", expectedPath, ex.Path)
	}
}

func TestGitExecution_Parse_Clone(t *testing.T) {
	workspace := config.Workspace{Name: "ws", Path: "/tmp/ws"}
	group := config.Group{Name: "group"}
	repo := config.Repository{Name: "repo", Remotes: []config.GitRemote{{Name: "origin", URL: "https://example.com/repo.git"}}}

	cmd := fakeCobraCommand("git")

	execs, err := GitExecution{}.Parse(cmd, []string{"clone"}, workspace, group, repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(execs) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(execs))
	}
	ex := execs[0]
	if ex.Command != "git" {
		t.Fatalf("expected command git, got %s", ex.Command)
	}
	// clone adds remote URL and repository name
	if len(ex.Args) != 3 || ex.Args[0] != "clone" || ex.Args[1] != "https://example.com/repo.git" || ex.Args[2] != "repo" {
		t.Fatalf("unexpected clone args: %#v", ex.Args)
	}
	// path should be group path for clone (parent directory of repo)
	expectedPath := GroupPath(workspace, group)
	if ex.Path != expectedPath {
		t.Fatalf("expected path %s, got %s", expectedPath, ex.Path)
	}
}

func TestGitExecution_Parse_Clone_WithCustomRepoPath(t *testing.T) {
	workspace := config.Workspace{Name: "ws", Path: "/tmp/ws"}
	group := config.Group{Name: "group"}
	// repository.Path includes nested directory and repo folder name
	repo := config.Repository{Name: "repo", Path: "nested/repo", Remotes: []config.GitRemote{{Name: "origin", URL: "https://example.com/repo.git"}}}

	cmd := fakeCobraCommand("git")

	execs, err := GitExecution{}.Parse(cmd, []string{"clone"}, workspace, group, repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(execs) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(execs))
	}
	ex := execs[0]
	// When custom path is provided, clone path should be parent directory (without trailing repo name)
	groupPath := GroupPath(workspace, group)
	expectedPath := groupPath + "/nested"
	if ex.Path != expectedPath {
		t.Fatalf("expected path %s, got %s", expectedPath, ex.Path)
	}
}

func TestGitExecution_Parse_RemoteCheck(t *testing.T) {
	workspace := config.Workspace{Name: "ws", Path: "/tmp/ws"}
	group := config.Group{Name: "group"}
	repo := config.Repository{Name: "repo", Remotes: []config.GitRemote{
		{Name: "origin", URL: "https://example.com/repo.git"},
		{Name: "upstream", URL: "https://example.com/upstream.git"},
	}}

	cmd := fakeCobraCommand("git")

	execs, err := GitExecution{}.Parse(cmd, []string{"remote", "check"}, workspace, group, repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// origin -> set-url, upstream -> remove + add
	if len(execs) != 3 {
		for i, ex := range execs { t.Logf("exec %d: %#v", i, ex) }
		t.Fatalf("expected 3 executions, got %d", len(execs))
	}
	if execs[0].Args[0] != "remote" || execs[0].Args[1] != "set-url" || execs[0].Args[2] != "origin" {
		t.Fatalf("unexpected first execution args: %#v", execs[0].Args)
	}
	if execs[1].Args[1] != "remove" || execs[1].Args[2] != "upstream" {
		t.Fatalf("unexpected second execution args: %#v", execs[1].Args)
	}
	if execs[2].Args[1] != "add" || execs[2].Args[2] != "upstream" {
		t.Fatalf("unexpected third execution args: %#v", execs[2].Args)
	}
	// All remote check executions run in repository path
	expectedPath := RepositoryPath(workspace, group, repo)
	for i, ex := range execs {
		if ex.Path != expectedPath {
			t.Fatalf("execution %d expected path %s, got %s", i, expectedPath, ex.Path)
		}
	}
}
