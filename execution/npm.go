package execution

import (
  "fmt"
  "os"
  "strings"

  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/config"
  "github.com/carrchang/handy-ci/util"
)

type NpmExecution struct {
}

func (s NpmExecution) CheckArgs(cmd *cobra.Command, args []string) error {
  return nil
}

func (s NpmExecution) Parse(
  cmd *cobra.Command, args []string,
  workspace config.Workspace, group config.Group, repository config.Repository) ([]Execution, error) {
  pkg, _ := cmd.Flags().GetString(util.HandyCiNpmFlagPackage)
  path := RepositoryPath(workspace, group, repository)

  var executions []Execution

  for _, npm := range repository.Npms {
    if pkg != "" && pkg != npm.Name {
      continue
    }

    executionPath := strings.TrimRight(
      fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", path, strings.Trim(npm.Path, string(os.PathSeparator))),
      string(os.PathSeparator))

    executions = append(executions, Execution{
      Command: cmd.Use,
      Path:    executionPath,
      Args:    args,
    })
  }

  return executions, nil
}
