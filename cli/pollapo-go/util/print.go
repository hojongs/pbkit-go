package util

import (
	"fmt"

	"github.com/fatih/color"
)

var Yellow = color.New(color.FgYellow).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()

func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(color.Output, format, a...)
}

func Print(a ...interface{}) (n int, err error) {
	return fmt.Print(a...)
}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func PrintfVerbose(name string, verbose bool, format string, a ...interface{}) (n int, err error) {
	if verbose {
		return Printf(fmt.Sprintf("VERBOSE[%s]: %s", name, format), a...)
	} else {
		return 0, nil
	}
}
