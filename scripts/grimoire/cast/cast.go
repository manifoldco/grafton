// Package cast provides utility funcs for magefiles.
package cast

import (
	"fmt"
	"strings"

	"github.com/google/shlex"
	"github.com/magefile/mage/sh"
)

// Sh formats the provided string and arguments, tokenizes the result using
// shell style rules, and then executes the result with mage's sh.RunWith.
func Sh(format string, a ...interface{}) error {
	envs, cmd, args, err := parse(format, a...)
	if err != nil {
		return err
	}

	return sh.RunWith(envs, cmd, args...)
}

// ShOutput formats the provided string and arguments, tokenizes the result using
// shell style rules, and then executes the result with mage's sh.OutputWith
func ShOutput(format string, a ...interface{}) (string, error) {
	envs, cmd, args, err := parse(format, a...)
	if err != nil {
		return "", err
	}

	return sh.OutputWith(envs, cmd, args...)
}

// PrepShOutput prepares an ShOut invocation with the provided prefix. when the returned
// func is executed, prefix and suffix are joined with a space, then passed in
// to Sh along with any values for a.
//
// If you only wish to pass values for a, use "" for suffix.
func PrepShOutput(prefix string) func(suffix string, a ...interface{}) (string, error) {
	return func(suffix string, a ...interface{}) (string, error) {
		format := prefix
		if suffix != "" {
			format = fmt.Sprintf("%s %s", format, suffix)
		}

		return ShOutput(format, a...)
	}
}

func parse(format string, a ...interface{}) (env map[string]string, cmd string, args []string, err error) {
	raw := fmt.Sprintf(format, a...)
	parts, err := shlex.Split(raw)
	if err != nil {
		return nil, "", nil, err
	}

	envs := make(map[string]string)
	for len(parts) > 0 && strings.Contains(parts[0], "=") {
		envParts := strings.SplitN(parts[0], "=", 2)
		k, v := envParts[0], envParts[1]
		envs[k] = v
		parts = parts[1:]
	}

	return envs, parts[0], parts[1:], nil
}
