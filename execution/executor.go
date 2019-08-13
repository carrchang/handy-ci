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
  "bufio"
  "context"
  "fmt"
  "os/exec"
  "strings"
  "time"

  "github.com/spf13/cobra"

  "github.com/carrchang/handy-ci/config"
  "github.com/carrchang/handy-ci/util"
)

func Execute(cmd *cobra.Command, args []string, executionParser Parser) {
  var err error

  args = util.ParseFlagsAndArgs(cmd.Flags(), args)

  err = executionParser.CheckArgs(cmd, args)

  if err != nil {
    fmt.Printf("Error:\n  %s\n", err.Error())
    cmd.Help()
    return
  }

  help, _ := cmd.Flags().GetBool("help")

  if help {
    cmd.Help()
    return
  }

  execInWorkspaces(cmd, args, executionParser)
}

func execInWorkspaces(cmd *cobra.Command, args []string, executionParser Parser) error {
  currentWorkspace, _ := cmd.Flags().GetString(util.HandyCiFlagWorkspace)

  for _, workspace := range util.Workspaces() {
    if currentWorkspace != "" {
      if workspace.Name == currentWorkspace {
        err := execInGroups(cmd, args, executionParser, workspace)

        if err != nil {
          return err
        }
      }
    } else {
      err := execInGroups(cmd, args, executionParser, workspace)

      if err != nil {
        return err
      }
    }
  }

  return nil
}

func execInGroups(cmd *cobra.Command, args []string, executionParser Parser, workspace config.Workspace) error {
  currentGroup, _ := cmd.Flags().GetString(util.HandyCiFlagGroup)

  for _, group := range workspace.Groups {
    if currentGroup != "" {
      if group.Name == currentGroup {
        err := execInRepositories(cmd, args, executionParser, workspace, group)

        if err != nil {
          return err
        }
      }
    } else {
      err := execInRepositories(cmd, args, executionParser, workspace, group)

      if err != nil {
        return err
      }
    }
  }

  return nil
}

func execInRepositories(
  cmd *cobra.Command, args []string, executionParser Parser, workspace config.Workspace, group config.Group) error {
  currentRepository, _ := cmd.Flags().GetString(util.HandyCiFlagRepository)
  toBeContinue, _ := cmd.Flags().GetBool(util.HandyCiFlagContinue)

  skippedRepositoriesInString, _ := cmd.Flags().GetString(util.HandyCiFlagSkip)

  var skippedRepositories []string

  for _, skippedRepository := range strings.Split(skippedRepositoriesInString, ",") {
    skippedRepositories = append(skippedRepositories, strings.Trim(skippedRepository, " "))
  }

  for _, repository := range group.Repositories {
    if util.ContainArgs(skippedRepositories, repository.Name) {
      continue
    }

    if currentRepository != "" {
      if repository.Name == currentRepository {
        i, err := execInRepository(cmd, args, executionParser, workspace, group, repository, toBeContinue)

        if err != nil && !toBeContinue {
          return err
        }

        if i > 0 {
          fmt.Println()
        }
      }
    } else {
      i, err := execInRepository(cmd, args, executionParser, workspace, group, repository, toBeContinue)

      if err != nil && !toBeContinue {
        return err
      }

      if i > 0 {
        fmt.Println()
      }
    }
  }

  return nil
}

func execInRepository(
  cmd *cobra.Command, args []string, executionParser Parser,
  workspace config.Workspace, group config.Group, repository config.Repository, toBeContinue bool) (int, error) {
  executions, err := executionParser.Parse(cmd, args, workspace, group, repository)

  for i, execution := range executions {
    if err != nil && !toBeContinue {
      util.Printf("%s", err.Error())
      return i, err
    }

    util.Printf("PATH: %s", execution.Path)
    util.Printf("CMD: %s %s", execution.Command, strings.Join(execution.Args, " "))

    if execution.Skip {
      continue
    }

    context, cancelFn := context.WithTimeout(context.Background(), time.Hour*24)
    defer cancelFn()

    cmd := exec.CommandContext(context, execution.Command, execution.Args...)
    cmd.Dir = execution.Path

    //cmd.SysProcAttr = &syscall.SysProcAttr{
    //	Setpgid: true,
    //}

    cmdReader, err := cmd.StdoutPipe()
    if err != nil && !toBeContinue {
      util.Printf("%s", err.Error())
      return i, err
    }

    done := make(chan struct{})

    scanner := bufio.NewScanner(cmdReader)
    go func() {
      for scanner.Scan() {
        util.Printf("%s", scanner.Text())
      }

      done <- struct{}{}
    }()

    err = cmd.Start()
    if err != nil && !toBeContinue {
      util.Printf("%s", err.Error())
      return i, err
    }

    <-done

    err = cmd.Wait()
    if err != nil && !toBeContinue {
      util.Printf("%s", err.Error())
      return i, err
    }

    if i < len(executions)-1 {
      fmt.Println()
    }
  }

  return len(executions), nil
}
