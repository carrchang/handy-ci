package execution

import (
  "fmt"
  "os"
  "strings"

  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/config"
  "github.com/carrchang/handy-ci/util"
)

type GitExecution struct {
}

func (s GitExecution) CheckArgs(command *cobra.Command, args []string) error {
  return nil
}

func (s GitExecution) Parse(
  command *cobra.Command, args []string,
  workspace config.Workspace, group config.Group, repository config.Repository) ([]Execution, error) {
  var path string
  var err error

  path = RepositoryPath(workspace, group, repository)

  if util.ContainArgs(args, "clone") {
    _, err = os.Stat(path)

    if os.IsNotExist(err) {
      err = nil
    }

    args = append(args, RepositoryRemoteURL(repository, "origin"))
    args = append(args, repository.Name)

    path = GroupPath(workspace, group)

    if repository.Path != "" {
      path = fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", path, repository.Path)
      path = strings.TrimSuffix(strings.TrimSuffix(path, repository.Name), string(os.PathSeparator))
    }

    _, err2 := os.Stat(path)
    if os.IsNotExist(err2) {
      os.MkdirAll(path, 0755)
    }
  }

  if util.ContainArgs(args, "remote") && util.ContainArgs(args, "check") {
    args = []string{}
    args = append(args, "remote")

    var executions []Execution
    var remote config.GitRemote
    var executionArgs []string

    for _, remote = range repository.Remotes {
      if remote.Name == "origin" {
        executionArgs = args
        executionArgs = append(executionArgs, "set-url")
        executionArgs = append(executionArgs, remote.Name)
        executionArgs = append(executionArgs, remote.URL)

        executions = append(executions, Execution{
          Command: command.Use,
          Path:    path,
          Args:    executionArgs,
        })
      } else {
        executionArgs = args
        executionArgs = append(executionArgs, "remove")
        executionArgs = append(executionArgs, remote.Name)

        executions = append(executions, Execution{
          Command: command.Use,
          Path:    path,
          Args:    executionArgs,
        })

        executionArgs = args
        executionArgs = append(executionArgs, "add")
        executionArgs = append(executionArgs, remote.Name)
        executionArgs = append(executionArgs, remote.URL)

        executions = append(executions, Execution{
          Command: command.Use,
          Path:    path,
          Args:    executionArgs,
        })
      }
    }

    return executions, nil
  }

  return []Execution{
    {
      Command: command.Use,
      Path:    path,
      Args:    args,
    },
  }, nil
}
