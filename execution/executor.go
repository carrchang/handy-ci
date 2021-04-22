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
  targetRepositoriesInString, _ := cmd.Flags().GetString(util.HandyCiFlagRepositories)

  var targetRepositories []string

  for _, targetRepository := range strings.Split(targetRepositoriesInString, ",") {
    targetRepositories = append(targetRepositories, strings.Trim(targetRepository, " "))
  }

  toBeContinue, _ := cmd.Flags().GetBool(util.HandyCiFlagContinue)
  fromRepository, _ := cmd.Flags().GetString(util.HandyCiFlagFrom)

  skippedRepositoriesInString, _ := cmd.Flags().GetString(util.HandyCiFlagSkip)

  var skippedRepositories []string

  for _, skippedRepository := range strings.Split(skippedRepositoriesInString, ",") {
    skippedRepositories = append(skippedRepositories, strings.Trim(skippedRepository, " "))
  }

  dryRun, _ := cmd.Flags().GetBool(util.HandyCiFlagDryRun)

  var resume bool

  for _, repository := range group.Repositories {
    if !resume && fromRepository != "" {
      if strings.EqualFold(repository.Name, fromRepository) {
        resume = true
      } else {
        continue
      }
    }

    if util.ContainArgs(skippedRepositories, repository.Name) {
      continue
    }

    if targetRepositoriesInString != "" {
      if util.ContainArgs(targetRepositories, repository.Name) {
        i, err := execInRepository(cmd, args, executionParser, workspace, group, repository, toBeContinue, dryRun)

        if err != nil && !toBeContinue {
          return err
        }

        if i > 0 {
          util.Println()
        }
      }
    } else {
      i, err := execInRepository(cmd, args, executionParser, workspace, group, repository, toBeContinue, dryRun)

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
  workspace config.Workspace, group config.Group, repository config.Repository, toBeContinue bool, dryRun bool) (int, error) {
  executions, err := executionParser.Parse(cmd, args, workspace, group, repository)

  for i, execution := range executions {
    if err != nil && !toBeContinue {
      util.Printf("%v\n", err.Error())
      return i, err
    }

    util.Printf("CMD: %s %s\n", execution.Command, strings.Join(execution.Args, " "))
    util.Printf("PATH: %s\n", execution.Path)

    if dryRun {
      continue
    }

    if execution.Skip {
      continue
    }

    util.Printf("%s\n", ">>>>>>")

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

    if (args[i] == "--"+util.HandyCiFlagRepositories || args[i] == "-"+util.HandyCiFlagRepositoriesShorthand) &&
      len(args) >= i+1 {
      flags.Set(util.HandyCiFlagRepositories, args[i+1])

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

    if args[i] == "--"+util.HandyCiFlagDryRun {
      flags.Set(util.HandyCiFlagDryRun, "true")

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

    if args[i] == "--"+util.HandyCiExecFlagNonStrict {
      flags.Set(util.HandyCiExecFlagNonStrict, "true")

      continue
    }

    cleanedArgs = append(cleanedArgs, args[i])
  }

  return cleanedArgs
}

type executionWriter struct {
}

func NewWriter() io.Writer {
  return executionWriter{}
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
