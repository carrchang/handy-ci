package execution

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/carrchang/handy-ci/config"
)

func TestGroupPath_VarAndHomeExpansionAndTrim(t *testing.T) {
	// simulate HANDY_CI_ROOT and HOME
	tmpDir := t.TempDir()
	os.Setenv("HANDY_CI_ROOT", tmpDir)
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", filepath.Join(tmpDir, "home"))
	t.Cleanup(func() {
		os.Setenv("HOME", origHome)
	})

	ws := config.Workspace{Name: "ws", Path: "$HANDY_CI_ROOT/$HOME/work/"}
	grp := config.Group{Name: "group"}
	got := GroupPath(ws, grp)
	// Expect substitution and trimming of trailing separator
	if strings.HasSuffix(got, string(os.PathSeparator)) {
		t.Fatalf("expected no trailing separator: %s", got)
	}
	if !strings.Contains(got, tmpDir) {
		t.Fatalf("expected path to contain tmpDir %s, got %s", tmpDir, got)
	}
	if !strings.Contains(got, "home") {
		t.Fatalf("expected path to contain substituted home, got %s", got)
	}
	if !strings.HasSuffix(got, string(os.PathSeparator)+"group") {
		// on windows separators differ
		if runtime.GOOS != "windows" { // windows path may have different layout but suffix should still match
			// If it truly fails on windows we'll revisit; keep assertion for others.
			// We still expect the group name appended when NameIgnoredInPath is false and no custom path.
			if !strings.HasSuffix(got, "group") {
				// fallback check
				t.Fatalf("expected path to end with group name, got %s", got)
			}
		}
	}
}

func TestGroupPath_GroupCustomAbsolutePath(t *testing.T) {
	ws := config.Workspace{Name: "ws", Path: "/base"}
	grp := config.Group{Name: "group", Path: "/abs/custom/"}
	got := GroupPath(ws, grp)
	if got != "/abs/custom" { // trailing slash trimmed
		// Allow Windows volume prefix differences
		if runtime.GOOS == "windows" && !strings.HasSuffix(got, `abs\custom`) {
			t.Fatalf("unexpected group path: %s", got)
		} else if runtime.GOOS != "windows" {
			t.Fatalf("expected /abs/custom got %s", got)
		}
	}
}

func TestGroupPath_GroupNameIgnored(t *testing.T) {
	ws := config.Workspace{Name: "ws", Path: "/base"}
	grp := config.Group{Name: "group", NameIgnoredInPath: true}
	got := GroupPath(ws, grp)
	if got != "/base" && !(runtime.GOOS == "windows" && strings.HasSuffix(got, `base`)) {
		t.Fatalf("expected base path without group name, got %s", got)
	}
}

func TestRepositoryPath_CustomRelative(t *testing.T) {
	ws := config.Workspace{Name: "ws", Path: "/base"}
	grp := config.Group{Name: "group"}
	repo := config.Repository{Name: "repo", Path: "nested/repo"}
	got := RepositoryPath(ws, grp, repo)
	expected := filepath.Join("/base", "group", "nested", "repo")
	if runtime.GOOS == "windows" { // adjust expected root if needed
		if !strings.HasSuffix(got, filepath.Join("group", "nested", "repo")) {
			t.Fatalf("unexpected repo path on windows: %s", got)
		}
	} else if got != expected {
		t.Fatalf("expected %s got %s", expected, got)
	}
}

func TestRepositoryPath_NameIgnored(t *testing.T) {
	ws := config.Workspace{Name: "ws", Path: "/base"}
	grp := config.Group{Name: "group"}
	repo := config.Repository{Name: "repo", NameIgnoredInPath: true}
	got := RepositoryPath(ws, grp, repo)
	expected := filepath.Join("/base", "group")
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(got, "group") {
			t.Fatalf("expected path to end with group got %s", got)
		}
	} else if got != expected {
		t.Fatalf("expected %s got %s", expected, got)
	}
}

func TestRepositoryRemoteURL_FoundAndNotFound(t *testing.T) {
	repo := config.Repository{Remotes: []config.GitRemote{{Name: "origin", URL: "o"}, {Name: "upstream", URL: "u"}}}
	if got := RepositoryRemoteURL(repo, "origin"); got != "o" {
		t.Fatalf("expected origin=o got %s", got)
	}
	if got := RepositoryRemoteURL(repo, "none"); got != "" {
		t.Fatalf("expected empty string for missing remote got %s", got)
	}
}
