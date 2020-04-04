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
  "bytes"
  "fmt"
  "io"
  "os/exec"
  "strings"
  "sync"

  "github.com/logrusorgru/aurora"
  "github.com/spf13/cobra"
  "github.com/spf13/pflag"

  "github.com/carrchang/handy-ci/config"
  "github.com/carrchang/handy-ci/util"
)

func Execute(cmd *cobra.Command, args []string, executionParser Parser) {
  var err error

  args = ParseFlagsAndArgs(cmd.Flags(), args)

  err = executionParser.CheckArgs(cmd, args)

  if err != nil {
    util.Printf("Error:\n  %v\n", err)
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

  for _, workspace := range Workspaces() {
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
  fromRepository, _ := cmd.Flags().GetString(util.HandyCiFlagFrom)

  skippedRepositoriesInString, _ := cmd.Flags().GetString(util.HandyCiFlagSkip)

  var skippedRepositories []string

  for _, skippedRepository := range strings.Split(skippedRepositoriesInString, ",") {
    skippedRepositories = append(skippedRepositories, strings.Trim(skippedRepository, " "))
  }

  var resume bool

  for _, repository := range group.Repositories {
    if util.ContainArgs(skippedRepositories, repository.Name) {
      continue
    }

    if fromRepository != "" {
      if !resume && repository.Name != fromRepository {
        continue
      } else {
        resume = true
      }
    }

    if currentRepository != "" {
      if repository.Name == currentRepository {
        i, err := execInRepository(cmd, args, executionParser, workspace, group, repository, toBeContinue)

        if err != nil && !toBeContinue {
          return err
        }

        if i > 0 {
          util.Println()
        }
      }
    } else {
      i, err := execInRepository(cmd, args, executionParser, workspace, group, repository, toBeContinue)

      if err != nil && !toBeContinue {
        return err
      }

      if i > 0 {
        util.Println()
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
      util.Printf("%v\n", err.Error())
      return i, err
    }

    util.Printf("PATH: %s\n", execution.Path)
    util.Printf("CMD: %s %s\n", execution.Command, strings.Join(execution.Args, " "))

    util.Printf("%s\n", ">>>>>>")

    if execution.Skip {
      continue
    }

    cmd := exec.Command(execution.Command, execution.Args...)
    cmd.Dir = execution.Path

    var stdoutBuf, stderrBuf bytes.Buffer
    stdoutIn, _ := cmd.StdoutPipe()
    stderrIn, _ := cmd.StderrPipe()

    handyOut := NewWriter()
    handyErr := NewWriter()

    var errStdout, errStderr error
    stdout := io.MultiWriter(handyOut, &stdoutBuf)
    stderr := io.MultiWriter(handyErr, &stderrBuf)
    err := cmd.Start()

    fmt.Print(aurora.Green("[Handy CI]"), " ")

    if err != nil && !toBeContinue {
      fmt.Printf("%v\n", err)
      util.Printf("%s\n", "<<<<<<")
      return i, err
    }

    var wg sync.WaitGroup
    wg.Add(1)

    var count int64
    var stdCount int64

    go func() {
      stdCount, errStdout = io.Copy(stdout, stdoutIn)
      count += stdCount

      wg.Done()
    }()

    stdCount, errStderr = io.Copy(stderr, stderrIn)
    count += stdCount
    wg.Wait()

    err = cmd.Wait()
    if err != nil && !toBeContinue {
      fmt.Printf("%v\n", err)
      util.Printf("%s\n", "<<<<<<")
      return i, err
    }

    if (errStdout != nil || errStderr != nil) && !toBeContinue {
      fmt.Printf("Failed to capture stdout or stderr, %v, %v\n", errStdout, errStderr)
      util.Printf("%s\n", "<<<<<<")
      return i, err
    }

    fmt.Printf("%s\n", "<<<<<<")

    if i < len(executions)-1 {
      fmt.Println()
    }
  }

  return len(executions), nil
}

func Workspaces() []config.Workspace {
  return config.HandyCiConfig.Workspaces
}

func Groups(workspaceName string) []config.Group {
  var groups []config.Group

  for _, workspace := range config.HandyCiConfig.Workspaces {
    if workspace.Name == workspaceName {
      groups = workspace.Groups
    }
  }

  return groups
}

func ParseFlagsAndArgs(flags *pflag.FlagSet, args []string) []string {
  var cleanedArgs []string

  for i := 0; i < len(args); i++ {
    if (args[i] == "--"+util.HandyCiFlagWorkspace || args[i] == "-"+util.HandyCiFlagWorkspaceShorthand) &&
      len(args) >= i+1 {
      flags.Set(util.HandyCiFlagWorkspace, args[i+1])

      i++

      continue
    }

    if (args[i] == "--"+util.HandyCiFlagGroup || args[i] == "-"+util.HandyCiFlagGroupShorthand) &&
      len(args) >= i+1 {
      flags.Set(util.HandyCiFlagGroup, args[i+1])

      i++

      continue
    }

    if (args[i] == "--"+util.HandyCiFlagRepository || args[i] == "-"+util.HandyCiFlagRepositoryShorthand) &&
      len(args) >= i+1 {
      flags.Set(util.HandyCiFlagRepository, args[i+1])

      i++

      continue
    }

    if args[i] == "--"+util.HandyCiFlagContinue || args[i] == "-"+util.HandyCiFlagContinueShorthand {
      flags.Set(util.HandyCiFlagContinue, "true")

      continue
    }

    if (args[i] == "--"+util.HandyCiFlagFrom || args[i] == "-"+util.HandyCiFlagFromShorthand) &&
      len(args) >= i+1 {
      flags.Set(util.HandyCiFlagFrom, args[i+1])

      i++

      continue
    }

    if args[i] == "--"+util.HandyCiFlagSkip && len(args) >= i+1 {
      flags.Set(util.HandyCiFlagSkip, args[i+1])

      i++

      continue
    }

    if args[i] == "--"+util.HandyCiFlagConfig && len(args) >= i+1 {
      flags.Set(util.HandyCiFlagConfig, args[i+1])

      i++

      continue
    }

    if args[i] == "--"+util.HandyCiFlagHelp {
      flags.Set(util.HandyCiFlagHelp, "true")

      continue
    }

    if (args[i] == "--"+util.HandyCiNpmFlagPackage) && len(args) >= i+1 {
      flags.Set(util.HandyCiNpmFlagPackage, args[i+1])

      i++

      continue
    }

    cleanedArgs = append(cleanedArgs, args[i])
  }

  return cleanedArgs
}

type executionWriter struct {
}

func NewWriter() io.Writer {
  return executionWriter{
  }
}

func (e executionWriter) Write(p []byte) (int, error) {
  output := fmt.Sprintf("%s", p)

  if strings.Contains(output, "\n") {
    outputArgs := strings.Split(output, "\n")

    var realOutputArgs []interface{}

    for i, outputArg := range outputArgs {
      realOutputArgs = append(realOutputArgs, outputArg)

      if i != len(outputArgs)-1 {
        realOutputArgs = append(realOutputArgs, "\n", aurora.Green("[Handy CI]"), " ")
      }
    }

    fmt.Print(realOutputArgs...)
  } else {
    fmt.Print(output)
  }

  return len(p), nil
}
