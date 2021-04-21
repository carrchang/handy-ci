package execution

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/config"
  "github.com/carrchang/handy-ci/util"
)

type MavenExecution struct {
}

func (s MavenExecution) CheckArgs(cmd *cobra.Command, args []string) error {
  return nil
}

func (s MavenExecution) Parse(
  cmd *cobra.Command, args []string,
  workspace config.Workspace, group config.Group, repository config.Repository) ([]Execution, error) {
  path := RepositoryPath(workspace, group, repository)

  _, err := os.Stat(fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", path, "pom.xml"))

  var skip bool
  if os.IsNotExist(err) {
    util.Printf("Repository %s in path %s is not an valid maven project, skipped.\n", repository.Name, path)
    skip = true
  }

  return []Execution{
    {
      Command: cmd.Use,
      Path:    path,
      Args:    args,
      Skip:    skip,
    },
  }, nil
}
