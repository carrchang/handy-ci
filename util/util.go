package util

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-homedir"
)

const HandyCiName = "handy-ci"

const HandyCiFlagWorkspace = "workspace"
const HandyCiFlagWorkspaceShorthand = "W"
const HandyCiFlagGroup = "group"
const HandyCiFlagGroupShorthand = "G"
const HandyCiFlagRepositories = "repositories"
const HandyCiFlagRepositoriesShorthand = "R"
const HandyCiFlagTags = "tags"
const HandyCiFlagFrom = "from"
const HandyCiFlagFromShorthand = "F"
const HandyCiFlagSkip = "skip"
const HandyCiFlagContinue = "continue"
const HandyCiFlagContinueShorthand = "C"
const HandyCiExecFlagNonStrict = "non-strict"
const HandyCiFlagConfig = "config"
const HandyCiFlagDryRun = "dry-run"
const HandyCiFlagHelp = "help"

func Printf(format string, a ...interface{}) (n int, err error) {
	output := fmt.Sprintf(format, a...)

	if output != "" {
		return fmt.Print(aurora.Green("[Handy CI]"), " ", output)
	} else {
		return fmt.Print()
	}
}

func Println(a ...interface{}) (n int, err error) {
	output := fmt.Sprint(a...)

	if output != "" {
		return fmt.Println(aurora.Green("[Handy CI]"), " ", output)
	} else {
		return fmt.Println()
	}
}

func ContainArgs(args []string, arg string) bool {
	for _, currentArg := range args {
		if currentArg == arg {
			return true
		}
	}

	return false
}

func Home() string {
	home, err := homedir.Dir()
	if err != nil {
		Println(err)
		os.Exit(1)
	}

	return home
}
