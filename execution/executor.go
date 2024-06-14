package execution

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/carrchang/handy-ci/config"
	"github.com/carrchang/handy-ci/util"
)

func Execute(command *cobra.Command, args []string, executionParser Parser) {
	cleanedArgs, err := ParseFlagsAndArgs(command.Flags(), args)

	if err != nil {
		fmt.Printf("\n%v\n\n", err)
		return
	}

	err = executionParser.CheckArgs(command, cleanedArgs)

	if err != nil {
		fmt.Printf("\n%v\n\n", err)
		return
	}

	help, _ := command.Flags().GetBool("help")

	if help {
		command.Help()
		return
	}

	execInWorkspaces(command, cleanedArgs, executionParser)
}

func execInWorkspaces(command *cobra.Command, args []string, executionParser Parser) error {
	currentWorkspace, _ := command.Flags().GetString(util.HandyCiFlagWorkspace)

	for _, workspace := range Workspaces() {
		if currentWorkspace != "" {
			if workspace.Name == currentWorkspace {
				err := execInGroups(command, args, executionParser, workspace)

				if err != nil {
					return err
				}
			}
		} else {
			err := execInGroups(command, args, executionParser, workspace)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func execInGroups(command *cobra.Command, args []string, executionParser Parser, workspace config.Workspace) error {
	currentGroup, _ := command.Flags().GetString(util.HandyCiFlagGroup)

	for _, group := range workspace.Groups {
		if currentGroup != "" {
			if group.Name == currentGroup {
				err := execInRepositories(command, args, executionParser, workspace, group)

				if err != nil {
					return err
				}
			}
		} else {
			err := execInRepositories(command, args, executionParser, workspace, group)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func execInRepositories(
	command *cobra.Command, args []string, executionParser Parser, workspace config.Workspace, group config.Group) error {
	var targetRepositories []string
	targetRepositoriesInString, _ := command.Flags().GetString(util.HandyCiFlagRepositories)
	for _, targetRepository := range strings.Split(targetRepositoriesInString, ",") {
		targetRepositories = append(targetRepositories, strings.Trim(targetRepository, " "))
	}

	var tagsAsArgument []string
	tagsAsArgumentInString, _ := command.Flags().GetString(util.HandyCiFlagTags)
	for _, tagAsArgumentInString := range strings.Split(tagsAsArgumentInString, ",") {
		tagsAsArgument = append(tagsAsArgument, strings.Trim(tagAsArgumentInString, " "))
	}

	toBeContinue, _ := command.Flags().GetBool(util.HandyCiFlagContinue)
	fromRepository, _ := command.Flags().GetString(util.HandyCiFlagFrom)

	var skippedRepositories []string
	skippedRepositoriesInString, _ := command.Flags().GetString(util.HandyCiFlagSkip)
	for _, skippedRepository := range strings.Split(skippedRepositoriesInString, ",") {
		skippedRepositories = append(skippedRepositories, strings.Trim(skippedRepository, " "))
	}

	dryRun, _ := command.Flags().GetBool(util.HandyCiFlagDryRun)

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

		if !repositoryTagsContainAllTagsAsArgument(repository, tagsAsArgument) {
			continue
		}

		if targetRepositoriesInString != "" {
			if util.ContainArgs(targetRepositories, repository.Name) {
				i, err := execInRepository(command, args, executionParser, workspace, group, repository, toBeContinue, dryRun)

				if err != nil && !toBeContinue {
					return err
				}

				if i > 0 {
					util.Println()
				}
			}
		} else {
			i, err := execInRepository(command, args, executionParser, workspace, group, repository, toBeContinue, dryRun)

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

func repositoryTagsContainAllTagsAsArgument(repository config.Repository, tagsAsArgument []string) bool {
	for _, tagAsArgument := range tagsAsArgument {
		if tagAsArgument != "" && !util.ContainArgs(repository.Tags, tagAsArgument) {
			return false
		}
	}

	return true
}

func execInRepository(
	command *cobra.Command, args []string, executionParser Parser,
	workspace config.Workspace, group config.Group, repository config.Repository, toBeContinue bool, dryRun bool) (int, error) {
	util.Printf("PATH: %s\n", repository.Name)
	executions, err := executionParser.Parse(command, args, workspace, group, repository)

	for i, execution := range executions {
		if err != nil && !toBeContinue {
			util.Printf("%v\n", err.Error())
			return i, err
		}

		util.Printf("SCRIPT: %s %s\n", execution.Command, strings.Join(execution.Args, " "))
		util.Printf("PATH: %s\n", execution.Path)

		if dryRun {
			continue
		}

		if execution.Skip {
			continue
		}

		util.Printf("%s\n", ">>>>>>")

		executionCommand := exec.Command(execution.Command, execution.Args...)
		executionCommand.Dir = execution.Path
		executionCommand.Stdin = os.Stdin
		executionCommand.Stdout = os.Stdout
		executionCommand.Stderr = os.Stderr

		err := executionCommand.Run()
		if err != nil && !toBeContinue {
			fmt.Printf("%v\n", err)
			util.Printf("%s\n", "<<<<<<")
			return i, err
		}

		util.Printf("%s\n", "<<<<<<")

		if i < len(executions)-1 {
			fmt.Println()
		}
	}

	return len(executions), nil
}

func ScriptDefinitions() []config.ScriptDefinition {
	return config.HandyCiConfig.ScriptDefinitions
}

func Workspaces() []config.Workspace {
	return config.HandyCiConfig.Workspaces
}

func ParseFlagsAndArgs(flags *pflag.FlagSet, args []string) ([]string, error) {
	var cleanedArgs []string

	for i := 0; i < len(args); i++ {
		if args[i] == "--"+util.HandyCiFlagWorkspace || args[i] == "-"+util.HandyCiFlagWorkspaceShorthand {
			arg, err := parseFlagAndArg(args, i, args[i], true)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagWorkspace, arg)

			i++

			continue
		}

		if args[i] == "--"+util.HandyCiFlagGroup || args[i] == "-"+util.HandyCiFlagGroupShorthand {
			arg, err := parseFlagAndArg(args, i, args[i], true)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagGroup, arg)

			i++

			continue
		}

		if args[i] == "--"+util.HandyCiFlagRepositories || args[i] == "-"+util.HandyCiFlagRepositoriesShorthand {
			arg, err := parseFlagAndArg(args, i, args[i], true)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagRepositories, arg)

			i++

			continue
		}

		if args[i] == "--"+util.HandyCiFlagTags {
			arg, err := parseFlagAndArg(args, i, args[i], true)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagTags, arg)

			i++

			continue
		}

		if args[i] == "--"+util.HandyCiFlagFrom || args[i] == "-"+util.HandyCiFlagFromShorthand {
			arg, err := parseFlagAndArg(args, i, args[i], true)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagFrom, arg)

			i++

			continue
		}

		if args[i] == "--"+util.HandyCiFlagSkip {
			arg, err := parseFlagAndArg(args, i, args[i], true)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagSkip, arg)

			continue
		}

		if args[i] == "--"+util.HandyCiFlagContinue || args[i] == "-"+util.HandyCiFlagContinueShorthand {
			arg, err := parseFlagAndArg(args, i, args[i], false)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagContinue, arg)

			continue
		}

		if args[i] == "--"+util.HandyCiExecFlagNonStrict {
			arg, err := parseFlagAndArg(args, i, args[i], false)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiExecFlagNonStrict, arg)

			continue
		}

		if args[i] == "--"+util.HandyCiFlagConfig {
			arg, err := parseFlagAndArg(args, i, args[i], true)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagConfig, arg)

			i++

			continue
		}

		if args[i] == "--"+util.HandyCiFlagDryRun {
			arg, err := parseFlagAndArg(args, i, args[i], false)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagDryRun, arg)

			continue
		}

		if args[i] == "--"+util.HandyCiFlagHelp {
			arg, err := parseFlagAndArg(args, i, args[i], false)

			if err != nil {
				return cleanedArgs, err
			}

			flags.Set(util.HandyCiFlagHelp, arg)

			continue
		}

		cleanedArgs = append(cleanedArgs, args[i])
	}

	return cleanedArgs, nil
}

func parseFlagAndArg(args []string, i int, flag string, withValue bool) (string, error) {
	if withValue {
		if len(args) == i+1 {
			return "", ParseError{
				"Value for flag " + flag + " is required, use \"handy-ci --help\" for more information.",
			}
		} else {
			return args[i+1], nil
		}
	}

	return "true", nil
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
