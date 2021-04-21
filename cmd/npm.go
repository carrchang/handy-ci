package cmd

import (
  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/execution"
  "github.com/carrchang/handy-ci/util"
)

var npmCmd = &cobra.Command{
  Use:                "npm",
  Short:              "Execute npm command",
  DisableFlagParsing: true,
  Run: func(cmd *cobra.Command, args []string) {
    execution.Execute(cmd, args, execution.NpmExecution{})
  },
}

func init() {
  rootCmd.AddCommand(npmCmd)

  npmCmd.PersistentFlags().SortFlags = false
  npmCmd.Flags().SortFlags = false

  npmCmd.PersistentFlags().String(
    util.HandyCiNpmFlagPackage, "", "Npm package name in repository")
}
