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
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	"github.com/carrchang/handy-ci/util"
)

func SetupRootCommand(rootCmd *cobra.Command) {
	cobra.AddTemplateFunc("wrappedFlagUsages", wrappedFlagUsages)

	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpTemplate(helpTemplate)

	rootCmd.SuggestionsMinimumDistance = 2

	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	rootCmd.PersistentFlags().StringP(
		util.HandyCiFlagWorkspace, util.HandyCiFlagWorkspaceShorthand, "", "Execute command in workspace")
	rootCmd.PersistentFlags().StringP(
		util.HandyCiFlagGroup, util.HandyCiFlagGroupShorthand, "", "Execute command in group")
	rootCmd.PersistentFlags().StringP(
		util.HandyCiFlagRepository, util.HandyCiFlagRepositoryShorthand, "", "Execute command in repository")
	rootCmd.PersistentFlags().BoolP(
		util.HandyCiFlagContinue, util.HandyCiFlagContinueShorthand, false, "Skip failed command and continue")
	rootCmd.PersistentFlags().String(util.HandyCiFlagSkip, "", "Skip execution in comma-delimited list of repositories")
	rootCmd.PersistentFlags().String(util.HandyCiFlagConfig, "", "Config file (default is $HOME/.handy-ci/config.yaml)")
	rootCmd.PersistentFlags().Bool(util.HandyCiFlagHelp, false, "Print usage")
	rootCmd.PersistentFlags().Lookup(util.HandyCiFlagHelp).Hidden = true
}

type winSize struct {
	Height uint16
	Width  uint16
	x      uint16
	y      uint16
}

func getWinSize(fd uintptr) (*winSize, error) {
	uws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	ws := &winSize{Height: uws.Row, Width: uws.Col, x: uws.Xpixel, y: uws.Ypixel}
	return ws, err
}

func wrappedFlagUsages(cmd *cobra.Command) string {
	width := 80
	if ws, err := getWinSize(0); err == nil {
		width = int(ws.Width)
	}

	return cmd.Flags().FlagUsagesWrapped(width - 1)
}

var usageTemplate = `
{{- if ne .Long ""}}{{ .Long | trim }}{{- else }}{{ .Short | trim }}{{- end}}

Usage:

  {{.CommandPath}}{{- if .HasAvailableSubCommands}} COMMAND{{- end}} [OPTIONS]

{{- if .HasAvailableSubCommands}}

Commands:
{{- range .Commands}}
{{- if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}
{{- end}}
{{- end}}
{{- end}}

Options:
{{ wrappedFlagUsages . | trimRightSpace}}

You can use original options of "git", "mvn", "npm" or any command line tools as additional options.

{{- if .HasAvailableSubCommands}}

Use "{{.CommandPath}} COMMAND --help" for more information about a command.
{{- end}}

`

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
