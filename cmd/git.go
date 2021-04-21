package cmd

import (
  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/execution"
)

var gitCmd = &cobra.Command{
  Use:                "git",
  Short:              "Execute Git command",
  DisableFlagParsing: true,
  Run: func(cmd *cobra.Command, args []string) {
    execution.Execute(cmd, args, execution.GitExecution{})
  },
}

func init() {
  rootCmd.AddCommand(gitCmd)

  gitCmd.PersistentFlags().SortFlags = false
  gitCmd.Flags().SortFlags = false
}
