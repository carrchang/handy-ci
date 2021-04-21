package cmd

import (
  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/execution"
)

var mavenCmd = &cobra.Command{
  Use:                "mvn",
  Short:              "Execute Apache Maven command",
  DisableFlagParsing: true,
  Run: func(cmd *cobra.Command, args []string) {
    execution.Execute(cmd, args, execution.MavenExecution{})
  },
}

func init() {
  rootCmd.AddCommand(mavenCmd)

  mavenCmd.PersistentFlags().SortFlags = false
  mavenCmd.Flags().SortFlags = false
}
