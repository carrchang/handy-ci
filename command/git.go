package command

import (
  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/execution"
)

var gitCommand = &cobra.Command{
  Use:                "git",
  Short:              "Execute Git command",
  DisableFlagParsing: true,
  Run: func(command *cobra.Command, args []string) {
    execution.Execute(command, args, execution.GitExecution{})
  },
}

func init() {
  rootCommand.AddCommand(gitCommand)

  gitCommand.PersistentFlags().SortFlags = false
  gitCommand.Flags().SortFlags = false
}
