package acceptance

import (
	"fmt"
	"strings"

	"github.com/manifoldco/torus-cli/promptui"
)

var (
	bold  = promptui.Styler(promptui.FGBold)
	faint = promptui.Styler(promptui.FGFaint)
)

type resultCode int

const (
	pass resultCode = iota
	fail
)

// LogLevel is the type for setting our informational log level, irrespective of
// test results.
type LogLevel string

// Log Levels the user can set.
const (
	LogOff     LogLevel = "off"
	LogInfo    LogLevel = "info"
	LogVerbose LogLevel = "verbose"
)

func (l LogLevel) show(ol LogLevel) bool {
	switch {
	case l == "off":
		return false
	case l == "info":
		return ol != LogVerbose
	default:
		return true
	}
}

var lvl LogLevel

// SetLogLevel configures the log level for the duration of an acceptance test
// run.
func SetLogLevel(level LogLevel) {
	lvl = level
}

// Infof prints the given format string and args at the appropriate indentation
// level, if the  logger is at info level or above.
func Infof(format string, args ...interface{}) {
	if lvl.show(LogInfo) {
		printIndented(fmt.Sprintf(format, args...))
	}
}

// Infoln prints the given args at the appropriate indentation level, if the
// logger is at info level or above.
func Infoln(args ...interface{}) {
	if lvl.show(LogInfo) {
		printIndented(fmt.Sprintln(args...))
	}
}

var indent = 0
var entered = false

func enter(msg string) {
	printIndented(msg)
	entered = true
	indent += 2
}

func exit() {
	if entered {
		fmt.Println()
		entered = false
	}

	indent -= 2
}

func result(name string, code resultCode) {
	if !entered {
		indent -= 2
		printIndented(faint(name))
		indent += 2
	}

	var msg string
	switch code {
	case pass:
		msg = promptui.IconGood
		success++
	case fail:
		msg = promptui.IconBad
		failures++
	}

	fmt.Println(" " + msg)
	entered = false
}

func printIndented(msg string) {
	if entered {
		fmt.Println()
		entered = false
	}

	prefix := strings.Repeat(" ", indent)
	parts := strings.Split(msg, "\n")
	for i, part := range parts {
		if i == len(parts)-1 && part == "" {
			continue
		}

		fmt.Print(prefix + part)

		if i != len(parts)-1 {
			fmt.Println()
		}
	}
}

func printSummary(fails, success int) {
	styler := promptui.Styler(promptui.FGGreen)
	if fails > 0 {
		styler = promptui.Styler(promptui.FGRed)
	}
	msg := fmt.Sprintf("%d features, %d failures", fails+success, fails)
	fmt.Println()
	enter(styler(msg))
	exit()
}
