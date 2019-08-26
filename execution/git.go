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
	"os"

	"github.com/spf13/cobra"

	"github.com/carrchang/handy-ci/config"
	"github.com/carrchang/handy-ci/util"
)

type GitExecution struct {
}

func (s GitExecution) CheckArgs(cmd *cobra.Command, args []string) error {
	return nil
}

func (s GitExecution) Parse(
	cmd *cobra.Command, args []string,
	workspace config.Workspace, group config.Group, repository config.Repository) ([]Execution, error) {
	var path string
	var err error

	path = RepositoryPath(workspace, group, repository)

	if util.ContainArgs(args, "clone") {
		_, err = os.Stat(path)

		if os.IsNotExist(err) {
			err = nil
		}

		path = GroupPath(workspace, group)
		args = append(args, RepositoryRemoteURL(repository, "origin"))
		args = append(args, repository.Name)

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
					Command: cmd.Use,
					Path:    path,
					Args:    executionArgs,
				})
			} else {
				executionArgs = args
				executionArgs = append(executionArgs, "remove")
				executionArgs = append(executionArgs, remote.Name)

				executions = append(executions, Execution{
					Command: cmd.Use,
					Path:    path,
					Args:    executionArgs,
				})

				executionArgs = args
				executionArgs = append(executionArgs, "add")
				executionArgs = append(executionArgs, remote.Name)
				executionArgs = append(executionArgs, remote.URL)

				executions = append(executions, Execution{
					Command: cmd.Use,
					Path:    path,
					Args:    executionArgs,
				})
			}
		}

		return executions, nil
	}

	return []Execution{
		{
			Command: cmd.Use,
			Path:    path,
			Args:    args,
		},
	}, nil
}
