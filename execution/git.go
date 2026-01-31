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
    baseArgs := []string{"remote"}

    var executions []Execution
    var remote config.GitRemote

    for _, remote = range repository.Remotes {
      if remote.Name == "origin" {
        executionArgs := make([]string, len(baseArgs))
        copy(executionArgs, baseArgs)
        executionArgs = append(executionArgs, "set-url")
        executionArgs = append(executionArgs, remote.Name)
        executionArgs = append(executionArgs, remote.URL)

        executions = append(executions, Execution{
          Command: command.Use,
          Path:    path,
          Args:    executionArgs,
        })
      } else {
        removeArgs := make([]string, len(baseArgs))
        copy(removeArgs, baseArgs)
        removeArgs = append(removeArgs, "remove")
        removeArgs = append(removeArgs, remote.Name)

        executions = append(executions, Execution{
          Command: command.Use,
          Path:    path,
          Args:    removeArgs,
        })

        addArgs := make([]string, len(baseArgs))
        copy(addArgs, baseArgs)
        addArgs = append(addArgs, "add")
        addArgs = append(addArgs, remote.Name)
        addArgs = append(addArgs, remote.URL)

        executions = append(executions, Execution{
          Command: command.Use,
          Path:    path,
          Args:    addArgs,
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
