package execution

import (
  "fmt"
  "os"
  "path/filepath"
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
  CheckArgs(command *cobra.Command, args []string) error

  Parse(
    command *cobra.Command, args []string,
    workspace config.Workspace, group config.Group, repository config.Repository) ([]Execution, error)
}

type ParseError struct {
  message string
}

func (e ParseError) Error() string {
  return fmt.Sprintf("%s", e.message)
}

func GroupPath(workspace config.Workspace, group config.Group) string {
  workspacePath := filepath.FromSlash(workspace.Path)

  workspacePath = strings.ReplaceAll(workspacePath, "$HANDY_CI_ROOT", os.Getenv("HANDY_CI_ROOT"))

  homeDir, _ := os.UserHomeDir()

  if homeDir != "" {
    workspacePath = strings.ReplaceAll(workspacePath, "$HOME", homeDir)
  }

  workspacePath = strings.TrimSuffix(workspacePath, string(os.PathSeparator))

  if len(group.Path) > 0 {
    if strings.HasPrefix(group.Path, string(os.PathSeparator)) {
      return strings.TrimSuffix(group.Path, string(os.PathSeparator))
    } else {
      return fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", workspacePath, group.Path)
    }
  }

  if group.NameIgnoredInPath {
    return workspacePath
  } else {
    return fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", workspacePath, group.Name)
  }
}

func RepositoryPath(workspace config.Workspace, group config.Group, repository config.Repository) string {
  groupPath := GroupPath(workspace, group)

  groupPath = strings.TrimSuffix(groupPath, string(os.PathSeparator))

  if len(repository.Path) > 0 {
    if strings.HasPrefix(repository.Path, string(os.PathSeparator)) {
      return strings.TrimSuffix(repository.Path, string(os.PathSeparator))
    } else {
      return fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", groupPath, repository.Path)
    }
  }

  if repository.NameIgnoredInPath {
    return groupPath
  } else {
    return fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", groupPath, repository.Name)
  }
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
