package cmd

import (
  "github.com/carrchang/handy-ci/util"
  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/execution"
)

var execCmd = &cobra.Command{
  Use:                "exec",
  Short:              "Execute any command",
  DisableFlagParsing: true,
  Run: func(cmd *cobra.Command, args []string) {
    execution.Execute(cmd, args, execution.ExecExecution{})
  },
}

func init() {
  rootCmd.AddCommand(execCmd)

  execCmd.PersistentFlags().SortFlags = false
  execCmd.Flags().SortFlags = false

  execCmd.PersistentFlags().Bool(util.HandyCiExecFlagNonStrict, false, "Try to execute undefined command")
}
