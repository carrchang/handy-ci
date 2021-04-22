package execution

import (
  "fmt"
  "github.com/carrchang/handy-ci/util"
  "github.com/spf13/cobra"
  "os"
  "strings"

  "github.com/carrchang/handy-ci/config"
)

type ExecExecution struct {
}

func (s ExecExecution) CheckArgs(command *cobra.Command, args []string) error {
  return nil
}

func (s ExecExecution) Parse(
  command *cobra.Command, args []string,
  workspace config.Workspace, group config.Group, repository config.Repository) ([]Execution, error) {
  repositoryPath := RepositoryPath(workspace, group, repository)

  var currentScript string
  var executionArgs []string

  nonStrict, _ := command.Flags().GetBool(util.HandyCiExecFlagNonStrict)

  for i := 0; i < len(args); i++ {
    if i == 0 {
      currentScript = args[i]
    } else {
      executionArgs = append(executionArgs, args[i])
    }
  }

  var executions []Execution

  if currentScript != "" {
    var matched bool

    if len(repository.Scripts) > 0 {
      for _, script := range repository.Scripts {
        if !nonStrict {
          var scriptDefined bool

          for _, scriptDefinition := range ScriptDefinitions() {
            if scriptDefinition.Name == script.Name {
              scriptDefined = true

              break
            }
          }

          if !scriptDefined {
            util.Printf("Script %s not defined", script.Name)

            return executions, nil
          }
        }

        if currentScript == script.Name {
          if len(script.Paths) > 0 {
            for _, path := range script.Paths {
              executionPath := strings.TrimRight(
                fmt.Sprintf("%s"+string(os.PathSeparator)+"%s", repositoryPath, strings.Trim(path, string(os.PathSeparator))),
                string(os.PathSeparator))

              executions = append(executions, Execution{
                Command: currentScript,
                Path:    executionPath,
                Args:    executionArgs,
              })
            }
          } else {
            executions = append(executions, Execution{
              Command: currentScript,
              Path:    repositoryPath,
              Args:    executionArgs,
            })
          }

          matched = true

          break
        }
      }
    }

    if nonStrict && !matched {
      executions = append(executions, Execution{
        Command: currentScript,
        Path:    repositoryPath,
        Args:    executionArgs,
      })
    }
  } else {
    if len(repository.Scripts) > 0 {
      currentScript = repository.Scripts[0].Name

      for _, script := range repository.Scripts {
        if script.Default {
          currentScript = script.Name
        }
      }

      for _, scriptDefinition := range ScriptDefinitions() {
        if scriptDefinition.Name == currentScript {
          var args = strings.Split(scriptDefinition.DefaultArgs, " ")

          for _, arg := range args {
            executionArgs = append(executionArgs, strings.Trim(arg, " "))
          }
        }
      }

      executions = append(executions, Execution{
        Command: currentScript,
        Path:    repositoryPath,
        Args:    executionArgs,
      })
    }
  }

  return executions, nil
}
