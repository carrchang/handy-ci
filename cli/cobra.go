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
package cli

import (
  "fmt"
  "strings"

  "github.com/pkg/errors"
  "github.com/spf13/cobra"
  "golang.org/x/sys/unix"

  "github.com/carrchang/handy-ci/util"
)

func SetupRootCommand(rootCmd *cobra.Command) {
  cobra.AddTemplateFunc("hasSubCommands", hasSubCommands)
  cobra.AddTemplateFunc("hasManagementSubCommands", hasManagementSubCommands)
  cobra.AddTemplateFunc("operationSubCommands", operationSubCommands)
  cobra.AddTemplateFunc("managementSubCommands", managementSubCommands)
  cobra.AddTemplateFunc("wrappedFlagUsages", wrappedFlagUsages)
  cobra.AddTemplateFunc("customizedUsage", customizedUsage)

  rootCmd.SetUsageTemplate(usageTemplate)
  rootCmd.SetHelpTemplate(helpTemplate)
  rootCmd.SetFlagErrorFunc(FlagErrorFunc)
  rootCmd.SetHelpCommand(helpCommand)
  rootCmd.SetVersionTemplate("Handy CI version {{.Version}}\n")

  rootCmd.SuggestionsMinimumDistance = 1

  rootCmd.PersistentFlags().StringP(
    util.HandyCiFlagWorkspace, util.HandyCiFlagWorkspaceShorthand, "", "Execute command in workspace")
  rootCmd.PersistentFlags().StringP(
    util.HandyCiFlagGroup, util.HandyCiFlagGroupShorthand, "", "Execute command in group")
  rootCmd.PersistentFlags().StringP(
    util.HandyCiFlagRepository, util.HandyCiFlagRepositoryShorthand, "", "Execute command in repository")
  rootCmd.PersistentFlags().BoolP(
    util.HandyCiFlagContinue, util.HandyCiFlagContinueShorthand, false, "Skip failed command and continue")
  rootCmd.PersistentFlags().String(util.HandyCiFlagConfig, "", "Config file (default is $HOME/.handy-ci/config.yaml)")
  rootCmd.PersistentFlags().Bool(util.HandyCiFlagHelp, false, "Print usage")
  rootCmd.PersistentFlags().Lookup(util.HandyCiFlagHelp).Hidden = true
  rootCmd.PersistentFlags().SortFlags = false
}

// FlagErrorFunc prints an error message which matches the format of the
// docker/docker/cli error messages
func FlagErrorFunc(cmd *cobra.Command, err error) error {
  if err == nil {
    return nil
  }

  usage := ""
  if cmd.HasSubCommands() {
    usage = cmd.UsageString()
  }
  return StatusError{
    Status:     fmt.Sprintf("%s\nSee '%s --help'.%s", err, cmd.CommandPath(), usage),
    StatusCode: 125,
  }
}

func customizedUsage(cmd *cobra.Command) string {
  commands := strings.Split(cmd.CommandPath(), " ")

  usage := strings.Trim(commands[0], " ")

  if len(commands) > 1 {
    usage += " [Global Options] " + commands[1]
  } else {
    usage += " [Options] COMMAND"
  }

  return usage
}

func hasSubCommands(cmd *cobra.Command) bool {
  return len(operationSubCommands(cmd)) > 0
}

func hasManagementSubCommands(cmd *cobra.Command) bool {
  return len(managementSubCommands(cmd)) > 0
}

func operationSubCommands(cmd *cobra.Command) []*cobra.Command {
  var cmds []*cobra.Command
  for _, sub := range cmd.Commands() {
    if sub.IsAvailableCommand() && !sub.HasSubCommands() {
      cmds = append(cmds, sub)
    }
  }
  return cmds
}

type WinSize struct {
  Height uint16
  Width  uint16
  x      uint16
  y      uint16
}

func GetWinSize(fd uintptr) (*WinSize, error) {
  uws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
  ws := &WinSize{Height: uws.Row, Width: uws.Col, x: uws.Xpixel, y: uws.Ypixel}
  return ws, err
}

func wrappedFlagUsages(cmd *cobra.Command) string {
  width := 80
  if ws, err := GetWinSize(0); err == nil {
    width = int(ws.Width)
  }
  return cmd.Flags().FlagUsagesWrapped(width - 1)
}

func managementSubCommands(cmd *cobra.Command) []*cobra.Command {
  var cmds []*cobra.Command
  for _, sub := range cmd.Commands() {
    if sub.IsAvailableCommand() && sub.HasSubCommands() {
      cmds = append(cmds, sub)
    }
  }
  return cmds
}

var helpCommand = &cobra.Command{
  Use:               "help [command]",
  Short:             "Help about the command",
  Hidden:            true,
  PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
  PersistentPostRun: func(cmd *cobra.Command, args []string) {},
  RunE: func(c *cobra.Command, args []string) error {
    cmd, args, e := c.Root().Find(args)
    if cmd == nil || e != nil || len(args) > 0 {
      return errors.Errorf("unknown help topic: %v", strings.Join(args, " "))
    }

    helpFunc := cmd.HelpFunc()
    helpFunc(cmd, args)
    return nil
  },
}

var usageTemplate = `{{ .Short | trim }}

Usage:
  {{customizedUsage .}}
{{if .HasAvailableSubCommands}}
Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{end}}{{if .HasAvailableLocalFlags}}
Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Options:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}

You can use original options of "git", "mvn", "npm" or any command line tools as additional options.{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} COMMAND --help" for more information about a command.{{end}}

`

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
