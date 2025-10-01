package execution

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/carrchang/handy-ci/config"
	"github.com/carrchang/handy-ci/util"
)

// fakeParser implements Parser for executor tests.
// It records received args and returns configured executions / error.
type fakeParser struct {
	parseErr    error
	executions  []Execution
	checkArgsErr error
}

func (f *fakeParser) CheckArgs(command *cobra.Command, args []string) error {
	return f.checkArgsErr
}

func (f *fakeParser) Parse(command *cobra.Command, args []string, ws config.Workspace, g config.Group, r config.Repository) ([]Execution, error) {
	return f.executions, f.parseErr
}

func TestParseFlagAndArg_WithValueMissing(t *testing.T) {
	_, err := parseFlagAndArg([]string{"--workspace"}, 0, "--workspace", true)
	if err == nil {
		t.Fatalf("expected error when flag value missing")
	}
}

func TestParseFlagsAndArgs_AllFlags(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	// define flags expected by ParseFlagsAndArgs
	flags.String(util.HandyCiFlagWorkspace, "", "")
	flags.String(util.HandyCiFlagGroup, "", "")
	flags.String(util.HandyCiFlagRepositories, "", "")
	flags.String(util.HandyCiFlagTags, "", "")
	flags.String(util.HandyCiFlagFrom, "", "")
	flags.String(util.HandyCiFlagSkip, "", "")
	flags.Bool(util.HandyCiFlagContinue, false, "")
	flags.String(util.HandyCiFlagConfig, "", "")
	flags.Bool(util.HandyCiFlagDryRun, false, "")
	flags.Bool(util.HandyCiFlagHelp, false, "")

	args := []string{
		"--workspace", "w1",
		"--group", "g1",
		"--repositories", "r1,r2",
		"--tags", "t1,t2",
		"--from", "r2",
		"--skip", "r3",
		"--continue",
		"--config", "cfg.yml",
		"--dry-run",
		"--help",
		"extra1", "extra2",
	}

	cleaned, err := ParseFlagsAndArgs(flags, args)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !reflect.DeepEqual(cleaned, []string{"extra1", "extra2"}) {
		t.Fatalf("unexpected cleaned args: %#v", cleaned)
	}
	// spot check a couple flags
	if v, _ := flags.GetString(util.HandyCiFlagWorkspace); v != "w1" { t.Fatalf("workspace not set: %s", v) }
	if v, _ := flags.GetBool(util.HandyCiFlagDryRun); !v { t.Fatalf("dry-run not true") }
}

func TestRepositoryTagsContainAllTagsAsArgument(t *testing.T) {
	repo := config.Repository{Tags: []string{"a", "b"}}
	if !repositoryTagsContainAllTagsAsArgument(repo, []string{"a"}) { t.Fatalf("expected true") }
	if repositoryTagsContainAllTagsAsArgument(repo, []string{"c"}) { t.Fatalf("expected false") }
	if repositoryTagsContainAllTagsAsArgument(repo, []string{"a", "c"}) { t.Fatalf("expected false") }
}

func TestExecInRepository_DryRunAndSkip(t *testing.T) {
	p := &fakeParser{executions: []Execution{{Command: "echo", Args: []string{"hi"}, Path: "./", Skip: true}}}
	cmd := &cobra.Command{Use: "test"}
	ws := config.Workspace{Name: "ws"}
	grp := config.Group{Name: "g"}
	repo := config.Repository{Name: "r"}
	count, err := execInRepository(cmd, nil, p, ws, grp, repo, false, true)
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if count != 1 { t.Fatalf("expected executions count 1 got %d", count) }
}

func TestExecInRepositories_Filters(t *testing.T) {
	p := &fakeParser{executions: []Execution{{Command: "echo", Args: []string{"hi"}, Path: "./"}}}
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String(util.HandyCiFlagRepositories, "", "")
	cmd.Flags().String(util.HandyCiFlagTags, "", "")
	cmd.Flags().String(util.HandyCiFlagFrom, "", "")
	cmd.Flags().String(util.HandyCiFlagSkip, "", "")
	cmd.Flags().Bool(util.HandyCiFlagContinue, false, "")
	cmd.Flags().Bool(util.HandyCiFlagDryRun, true, "")

	ws := config.Workspace{Name: "ws"}
	grp := config.Group{Name: "g", Repositories: []config.Repository{
		{Name: "a", Tags: []string{"t1"}},
		{Name: "b", Tags: []string{"t1", "t2"}},
	}}

	// set filters: only b by name and tag t2
	cmd.Flags().Set(util.HandyCiFlagRepositories, "b")
	cmd.Flags().Set(util.HandyCiFlagTags, "t2")

	if err := execInRepositories(cmd, nil, p, ws, grp); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestExecInRepositories_ResumeFrom(t *testing.T) {
	p := &fakeParser{executions: []Execution{{Command: "echo", Args: []string{"hi"}, Path: "./"}}}
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String(util.HandyCiFlagRepositories, "", "")
	cmd.Flags().String(util.HandyCiFlagTags, "", "")
	cmd.Flags().String(util.HandyCiFlagFrom, "", "")
	cmd.Flags().String(util.HandyCiFlagSkip, "", "")
	cmd.Flags().Bool(util.HandyCiFlagContinue, false, "")
	cmd.Flags().Bool(util.HandyCiFlagDryRun, true, "")

	ws := config.Workspace{Name: "ws"}
	grp := config.Group{Name: "g", Repositories: []config.Repository{{Name: "a"}, {Name: "b"}, {Name: "c"}}}
	cmd.Flags().Set(util.HandyCiFlagFrom, "b")

	if err := execInRepositories(cmd, nil, p, ws, grp); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestExecInGroups_WorkspaceFiltering(t *testing.T) {
	p := &fakeParser{executions: []Execution{{Command: "echo", Path: "./"}}}
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String(util.HandyCiFlagGroup, "", "")
	cmd.Flags().String(util.HandyCiFlagRepositories, "", "")
	cmd.Flags().String(util.HandyCiFlagTags, "", "")
	cmd.Flags().String(util.HandyCiFlagFrom, "", "")
	cmd.Flags().String(util.HandyCiFlagSkip, "", "")
	cmd.Flags().Bool(util.HandyCiFlagContinue, false, "")
	cmd.Flags().Bool(util.HandyCiFlagDryRun, true, "")

	ws := config.Workspace{Name: "ws", Groups: []config.Group{{Name: "g1"}, {Name: "g2"}}}

	// Execute without filter
	if err := execInGroups(cmd, nil, p, ws); err != nil { t.Fatalf("unexpected err: %v", err) }

	// Filter to g2
	cmd.Flags().Set(util.HandyCiFlagGroup, "g2")
	if err := execInGroups(cmd, nil, p, ws); err != nil { t.Fatalf("unexpected err: %v", err) }
}

func TestExecInWorkspaces_Filter(t *testing.T) {
	// Patch global config for this test
	old := config.HandyCiConfig
	config.HandyCiConfig = &config.Config{Workspaces: []config.Workspace{{Name: "w1"}, {Name: "w2"}}}
	defer func() { config.HandyCiConfig = old }()

	p := &fakeParser{executions: []Execution{{Command: "echo", Path: "./"}}}
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String(util.HandyCiFlagWorkspace, "", "")
	cmd.Flags().String(util.HandyCiFlagGroup, "", "")
	cmd.Flags().String(util.HandyCiFlagRepositories, "", "")
	cmd.Flags().String(util.HandyCiFlagTags, "", "")
	cmd.Flags().String(util.HandyCiFlagFrom, "", "")
	cmd.Flags().String(util.HandyCiFlagSkip, "", "")
	cmd.Flags().Bool(util.HandyCiFlagContinue, false, "")
	cmd.Flags().Bool(util.HandyCiFlagDryRun, true, "")

	if err := execInWorkspaces(cmd, nil, p); err != nil { t.Fatalf("unexpected err: %v", err) }
	cmd.Flags().Set(util.HandyCiFlagWorkspace, "w2")
	if err := execInWorkspaces(cmd, nil, p); err != nil { t.Fatalf("unexpected err: %v", err) }
}

func TestExecute_CheckArgsError(t *testing.T) {
	p := &fakeParser{checkArgsErr: errors.New("bad args")}
	cmd := &cobra.Command{Use: "test"}
	Execute(cmd, []string{}, p)
	// No panic expected
}

func TestExecute_Help(t *testing.T) {
	p := &fakeParser{}
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Bool(util.HandyCiFlagHelp, false, "")
	Execute(cmd, []string{"--help"}, p)
}

func TestNewWriter(t *testing.T) {
	w := NewWriter()
	if w == nil { t.Fatalf("expected writer") }
	// simple Write call (newline and non-newline)
	fmt.Fprintf(w, "hello")
	fmt.Fprintf(w, "world\nagain")
}
