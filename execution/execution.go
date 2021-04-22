package execution

import (
  "fmt"
  "os"
  "strings"

  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/config"
)

type Execution struct {
  Command string
  Path    string
  Args    []string
  Skip    bool
}

type Parser interface {
  CheckArgs(cmd *cobra.Command, args []string) error

  Parse(
    cmd *cobra.Command, args []string,
    workspace config.Workspace, group config.Group, repository config.Repository) ([]Execution, error)
}

type ParseError struct {
  message string
}

func (e ParseError) Error() string {
  return fmt.Sprintf("%s", e.message)
}

func GroupPath(workspace config.Workspace, group config.Group) string {
  rootPath := workspace.Root

  rootPath = strings.TrimSuffix(rootPath, string(os.PathSeparator))

  if group.PathIgnored {
    return rootPath
  }

  return fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", rootPath, group.Name)
}

func RepositoryPath(workspace config.Workspace, group config.Group, repository config.Repository) string {
  groupPath := GroupPath(workspace, group)

  groupPath = strings.TrimSuffix(groupPath, string(os.PathSeparator))

  if repository.PathIgnored {
    return groupPath
  }

  return fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", groupPath, repository.Name)
}

func RepositoryRemoteURL(repository config.Repository, remoteName string) string {
  var remoteUrl string

  for _, remote := range repository.Remotes {
    if remote.Name == remoteName {
      remoteUrl = remote.URL
    }
  }

  return remoteUrl
}
