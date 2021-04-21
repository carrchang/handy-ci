package execution

import (
	"github.com/spf13/cobra"

	"github.com/carrchang/handy-ci/config"
)

type ExecExecution struct {
}

func (s ExecExecution) CheckArgs(cmd *cobra.Command, args []string) error {
	return nil
}

func (s ExecExecution) Parse(
	cmd *cobra.Command, args []string,
	workspace config.Workspace, group config.Group, repository config.Repository) ([]Execution, error) {
	path := RepositoryPath(workspace, group, repository)

	var command string
	var executionArgs []string

	for i := 0; i < len(args); i++ {
		if i == 0 {
			command = args[i]
		} else {
			executionArgs = append(executionArgs, args[i])
		}
	}

	return []Execution{
		{
			Command: command,
			Path:    path,
			Args:    executionArgs,
		},
	}, nil
}
