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
  "os"

  "github.com/logrusorgru/aurora"
  "github.com/mitchellh/go-homedir"
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
const HandyCiFlagFrom = "from"
const HandyCiFlagFromShorthand = "F"
const HandyCiFlagSkip = "skip"
const HandyCiFlagConfig = "config"
const HandyCiFlagHelp = "help"
const HandyCiNpmFlagPackage = "pkg"

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
  var currentArg string
  for _, currentArg = range args {
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
