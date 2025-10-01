package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/carrchang/handy-ci/util"
)

func TestRootCommand_Use(t *testing.T) {
	if rootCommand.Use != util.HandyCiName {
				// root command name should match constant
		t.Fatalf("expected root command use %s, got %s", util.HandyCiName, rootCommand.Use)
	}
}

func TestRootCommand_HelpContainsCustomUsage(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCommand.SetOut(buf)
	rootCommand.SetErr(buf)
	rootCommand.SetArgs([]string{"--help"})
	_ = rootCommand.Execute()
	out := buf.String()
	if !strings.Contains(out, "Options can be in front") {
		t.Fatalf("expected custom usage string in help output, got: %s", out)
	}
}

func TestRootCommand_SubcommandsRegistered(t *testing.T) {
	var foundExec, foundGit bool
	for _, c := range rootCommand.Commands() {
		if c == execCommand { foundExec = true }
		if c == gitCommand { foundGit = true }
	}
	if !foundExec || !foundGit {
		t.Fatalf("expected exec (%v) and git (%v) to be registered", foundExec, foundGit)
	}
}
