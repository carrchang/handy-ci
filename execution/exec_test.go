package execution

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/carrchang/handy-ci/config"
	"github.com/carrchang/handy-ci/util"
)

// helper to build an exec cobra command with needed flags
func newExecCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "exec"}
	cmd.Flags().Bool(util.HandyCiExecFlagNonStrict, false, "") // define non-strict flag for tests
	return cmd
}

func TestExecExecution_CheckArgs_NoArgs_NoPanic(t *testing.T) {
	config.HandyCiConfig = &config.Config{ // empty definitions are fine; no args triggers early return
		ScriptDefinitions: []config.ScriptDefinition{},
	}

	e := ExecExecution{}
	cmd := newExecCommand()

	// ensure no panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic, got %v", r)
		}
	}()

	if err := e.CheckArgs(cmd, []string{}); err != nil {
		// should be nil because empty args allowed now
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestExecExecution_CheckArgs_UnknownStrict_Error(t *testing.T) {
	config.HandyCiConfig = &config.Config{ScriptDefinitions: []config.ScriptDefinition{{Name: "mvn"}}}
	cmd := newExecCommand()
	err := ExecExecution{}.CheckArgs(cmd, []string{"unknown"})
	if err == nil {
		// expecting an error
		t.Fatalf("expected error for unknown script in strict mode")
	}
}

func TestExecExecution_CheckArgs_KnownStrict_OK(t *testing.T) {
	config.HandyCiConfig = &config.Config{ScriptDefinitions: []config.ScriptDefinition{{Name: "mvn"}}}
	cmd := newExecCommand()
	err := ExecExecution{}.CheckArgs(cmd, []string{"mvn"})
	if err != nil {
		// should not error when script is defined
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestExecExecution_CheckArgs_UnknownNonStrict_OK(t *testing.T) {
	config.HandyCiConfig = &config.Config{ScriptDefinitions: []config.ScriptDefinition{{Name: "mvn"}}}
	cmd := newExecCommand()
	cmd.Flags().Set(util.HandyCiExecFlagNonStrict, "true")
	err := ExecExecution{}.CheckArgs(cmd, []string{"other"})
	if err != nil {
		t.Fatalf("expected nil error in non-strict mode, got %v", err)
	}
}

func TestExecExecution_Parse_DefaultScriptSelection(t *testing.T) {
	config.HandyCiConfig = &config.Config{ScriptDefinitions: []config.ScriptDefinition{{Name: "npm", DefaultArgs: "outdated"}, {Name: "mvn", DefaultArgs: "clean install"}}}

	workspace := config.Workspace{Name: "ws", Path: "/root"}
	group := config.Group{Name: "grp"}
	repo := config.Repository{Name: "repo", Scripts: []config.Script{{Name: "npm"}, {Name: "mvn", Default: true}}}

	cmd := newExecCommand()

	executions, err := ExecExecution{}.Parse(cmd, []string{}, workspace, group, repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(executions) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(executions))
	}
	exe := executions[0]
	if exe.Command != "mvn" { // default:true should override first
		t.Fatalf("expected command mvn, got %s", exe.Command)
	}
	// args should come from ScriptDefinitions default args
	if len(exe.Args) != 2 || exe.Args[0] != "clean" || exe.Args[1] != "install" {
		t.Fatalf("unexpected args: %#v", exe.Args)
	}
}

func TestExecExecution_Parse_NonStrictUnknownScript(t *testing.T) {
	config.HandyCiConfig = &config.Config{ScriptDefinitions: []config.ScriptDefinition{{Name: "mvn"}}}
	workspace := config.Workspace{Name: "ws", Path: "/root"}
	group := config.Group{Name: "grp"}
	repo := config.Repository{Name: "repo"}
	cmd := newExecCommand()
	cmd.Flags().Set(util.HandyCiExecFlagNonStrict, "true")

	executions, err := ExecExecution{}.Parse(cmd, []string{"custom", "arg1"}, workspace, group, repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(executions) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(executions))
	}
	if executions[0].Command != "custom" || len(executions[0].Args) != 1 || executions[0].Args[0] != "arg1" {
		t.Fatalf("unexpected execution: %+v", executions[0])
	}
}
