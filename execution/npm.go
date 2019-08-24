/*
Copyright © 2019 Carr Chang <carr.z.chang@live.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
  path := util.RepositoryPath(workspace, group, repository)

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
