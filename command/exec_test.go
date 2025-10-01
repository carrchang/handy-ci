package command

import (
	"testing"

	"github.com/carrchang/handy-ci/util"
)

func TestExecCommand_HasNonStrictFlag(t *testing.T) {
	flag := execCommand.Flag(util.HandyCiExecFlagNonStrict)
	if flag == nil {
				// flag should be defined in exec.go init
		t.Fatalf("expected non-strict flag to exist on exec command")
	}
}
