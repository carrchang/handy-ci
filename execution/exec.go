package execution

import (
  "fmt"
  "github.com/carrchang/handy-ci/util"
  "github.com/spf13/cobra"
  "os"
  "strings"

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
  repositoryPath := RepositoryPath(workspace, group, repository)

  var command string
  var executionArgs []string

  nonStrict, _ := cmd.Flags().GetBool(util.HandyCiExecFlagNonStrict)

  for i := 0; i < len(args); i++ {
    if i == 0 {
      command = args[i]
    } else {
      executionArgs = append(executionArgs, args[i])
    }
  }

  var executions []Execution

  if command != "" {
    var matched bool

    if len(repository.Cmds) > 0 {
      for _, cmd := range repository.Cmds {
        if command == cmd.Name {
          if len(cmd.Paths) > 0 {
            for _, path := range cmd.Paths {
              executionPath := strings.TrimRight(
                fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", repositoryPath, strings.Trim(path, string(os.PathSeparator))),
                string(os.PathSeparator))

              executions = append(executions, Execution{
                Command: command,
                Path:    executionPath,
                Args:    executionArgs,
              })
            }
          } else {
            executions = append(executions, Execution{
              Command: command,
              Path:    repositoryPath,
              Args:    executionArgs,
            })
          }

          matched = true

          break
        }
      }
    }

    if nonStrict && !matched {
      executions = append(executions, Execution{
        Command: command,
        Path:    repositoryPath,
        Args:    executionArgs,
      })
    }
  }

  return executions, nil
}
