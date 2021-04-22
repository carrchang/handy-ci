package command

import (
  "github.com/carrchang/handy-ci/util"
  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/execution"
)

var execCommand = &cobra.Command{
  Use:                "exec",
  Short:              "Execute any command",
  DisableFlagParsing: true,
  Run: func(command *cobra.Command, args []string) {
    execution.Execute(command, args, execution.ExecExecution{})
  },
}

func init() {
  rootCommand.AddCommand(execCommand)

  execCommand.PersistentFlags().SortFlags = false
  execCommand.Flags().SortFlags = false

  execCommand.PersistentFlags().Bool(util.HandyCiExecFlagNonStrict, false, "Try to execute undefined command")
}
