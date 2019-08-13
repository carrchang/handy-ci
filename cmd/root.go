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
package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"

  "github.com/carrchang/handy-ci/config"
  "github.com/carrchang/handy-ci/util"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:                util.HandyCiName,
  Short:              "Handy CI is a tool for managing and building multi-repository source code on developer host.",
  DisableFlagParsing: true,
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
{{.Flags.FlagUsages | trimTrailingWhitespaces}}

Options can be in front of, behind, or on both sides of the command.

Original options of any command can be as additional options, and be in behind of the command.

{{- if .HasAvailableSubCommands}}

Use "{{.CommandPath}} COMMAND --help" for more information about a command.
{{- end}}

`

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

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

  configFlagUsage := "Config file (default is " + util.Home() +
    string(os.PathSeparator) + ".handy-ci" + string(os.PathSeparator) + "config.yaml)"

  rootCmd.PersistentFlags().String(util.HandyCiFlagConfig, "", configFlagUsage)
  rootCmd.PersistentFlags().Bool(util.HandyCiFlagHelp, false, "Print usage")
  rootCmd.PersistentFlags().Lookup(util.HandyCiFlagHelp).Hidden = true
}

func initConfig() {
  if cfgFile != "" {
    viper.SetConfigFile(cfgFile)
  } else {
    viper.AddConfigPath(util.Home() + string(os.PathSeparator) + "." + util.HandyCiName)
    viper.SetConfigName(util.HandyCiFlagConfig)
  }

  viper.AutomaticEnv()

  if err := viper.ReadInConfig(); err != nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
    fmt.Println(err)
  }

  config.Initialize()
}
