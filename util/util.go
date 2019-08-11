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
package util

import (
  "fmt"
  "strings"

  "github.com/spf13/pflag"

  "github.com/carrchang/handy-ci/config"
)

const HandyCiName = "handy-ci"

const HandyCiFlagWorkspace = "workspace"
const HandyCiFlagWorkspaceShorthand = "W"
const HandyCiFlagGroup = "group"
const HandyCiFlagGroupShorthand = "G"
const HandyCiFlagRepository = "repository"
const HandyCiFlagRepositoryShorthand = "R"
const HandyCiFlagContinue = "continue"
const HandyCiFlagContinueShorthand = "C"
const HandyCiFlagConfig = "config"
const HandyCiFlagHelp = "help"

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

func GroupPath(workspace config.Workspace, group config.Group) string {
  rootPath := workspace.Root

  rootPath = strings.TrimSuffix(rootPath, "/")

  return fmt.Sprintf("%s/%s", rootPath, group.Name)
}

func RepositoryPath(workspace config.Workspace, group config.Group, repository config.Repository) string {
  rootPath := GroupPath(workspace, group)

  rootPath = strings.TrimSuffix(rootPath, "/")

  return fmt.Sprintf("%s/%s", rootPath, repository.Name)
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

func Printf(format string, args ...interface{}) {
  message := fmt.Sprintf(format, args...)

  if message != "" {
    fmt.Println("[Handy CI] " + message)
  }
}

func ParseFlags(flags *pflag.FlagSet, args []string) {
  for i, arg := range args {
    if arg == "--"+HandyCiFlagConfig && len(args) >= i+1 {
      flags.Set(HandyCiFlagConfig, args[i+1])
    }

    if (arg == "--"+HandyCiFlagWorkspace || arg == "-"+HandyCiFlagWorkspaceShorthand) && len(args) >= i+1 {
      flags.Set(HandyCiFlagWorkspace, args[i+1])
    }

    if (arg == "--"+HandyCiFlagGroup || arg == "-"+HandyCiFlagGroupShorthand) && len(args) >= i+1 {
      flags.Set(HandyCiFlagGroup, args[i+1])
    }

    if (arg == "--"+HandyCiFlagRepository || arg == "-"+HandyCiFlagRepositoryShorthand) && len(args) >= i+1 {
      flags.Set(HandyCiFlagRepository, args[i+1])
    }

    if arg == "--"+HandyCiFlagContinue || arg == "-"+HandyCiFlagContinueShorthand {
      flags.Set(HandyCiFlagContinue, "true")
    }

    if arg == "--"+HandyCiFlagHelp {
      flags.Set(HandyCiFlagHelp, "true")
    }
  }
}

func CleanUpArgs(args []string) []string {
  var cleanedArgs []string

  for i := 0; i < len(args); i++ {
    if args[i] == "--"+HandyCiFlagConfig {
      i++
      continue
    }

    if args[i] == "--"+HandyCiFlagWorkspace || args[i] == "-"+HandyCiFlagWorkspaceShorthand {
      i++
      continue
    }

    if args[i] == "--"+HandyCiFlagGroup || args[i] == "-"+HandyCiFlagGroupShorthand {
      i++
      continue
    }

    if args[i] == "--"+HandyCiFlagRepository || args[i] == "-"+HandyCiFlagRepositoryShorthand {
      i++
      continue
    }

    if args[i] == "--"+HandyCiFlagContinue || args[i] == "-"+HandyCiFlagContinueShorthand {
      continue
    }

    if args[i] == "--"+HandyCiFlagHelp {
      continue
    }

    cleanedArgs = append(cleanedArgs, args[i])
  }

  return cleanedArgs
}
