package command

import (
  "os"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"

  "github.com/carrchang/handy-ci/config"
  "github.com/carrchang/handy-ci/util"
)

var cfgFile string

// rootCommand represents the base command when called without any subcommands
var rootCommand = &cobra.Command{
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
// This is called by main.main(). It only needs to happen once to the rootCommand.
func Execute() {
  if err := rootCommand.Execute(); err != nil {
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  rootCommand.SetUsageTemplate(usageTemplate)
  rootCommand.SetHelpTemplate(helpTemplate)

  rootCommand.SuggestionsMinimumDistance = 2

  rootCommand.PersistentFlags().SortFlags = false
  rootCommand.Flags().SortFlags = false

  rootCommand.PersistentFlags().StringP(
    util.HandyCiFlagWorkspace, util.HandyCiFlagWorkspaceShorthand, "", "Execute command in workspace")
  rootCommand.PersistentFlags().StringP(
    util.HandyCiFlagGroup, util.HandyCiFlagGroupShorthand, "", "Execute command in group")
  rootCommand.PersistentFlags().StringP(
    util.HandyCiFlagRepositories, util.HandyCiFlagRepositoriesShorthand,
    "", "Execute command in comma-delimited list of repositories")
  rootCommand.PersistentFlags().BoolP(
    util.HandyCiFlagContinue, util.HandyCiFlagContinueShorthand, false, "Skip failed command and continue")
  rootCommand.PersistentFlags().Bool(util.HandyCiFlagDryRun, false, "Only print the command and execution path")
  rootCommand.PersistentFlags().String(util.HandyCiFlagSkip, "", "Skip execution in comma-delimited list of repositories")
  rootCommand.PersistentFlags().StringP(
    util.HandyCiFlagFrom, util.HandyCiFlagFromShorthand, "", "Execute command from repository to end")

  configFlagUsage := "Config file (default is " + util.Home() +
    string(os.PathSeparator) + ".handy-ci" + string(os.PathSeparator) + "config.yaml)"

  rootCommand.PersistentFlags().String(util.HandyCiFlagConfig, "", configFlagUsage)
  rootCommand.PersistentFlags().Bool(util.HandyCiFlagHelp, false, "Print usage")
  rootCommand.PersistentFlags().Lookup(util.HandyCiFlagHelp).Hidden = true
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
    util.Println("Using config file:", viper.ConfigFileUsed())
    util.Println(err)
  }

  config.Initialize()
}
