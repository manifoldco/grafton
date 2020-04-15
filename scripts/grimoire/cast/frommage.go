package cast

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// These funcs are copied from mage, but updated so we can set the working
// directory

func runWith(dir string, env map[string]string, cmd string, args ...string) error {
	var output io.Writer
	if mg.Verbose() {
		output = os.Stdout
	}
	_, err := execute(dir, env, output, os.Stderr, cmd, args...)
	return err
}

func execute(dir string, env map[string]string, stdout, stderr io.Writer, cmd string, args ...string) (ran bool, err error) {
	expand := func(s string) string {
		s2, ok := env[s]
		if ok {
			return s2
		}
		return os.Getenv(s)
	}
	cmd = os.Expand(cmd, expand)
	for i := range args {
		args[i] = os.Expand(args[i], expand)
	}
	ran, code, err := run(dir, env, stdout, stderr, cmd, args...)
	if err == nil {
		return true, nil
	}
	if ran {
		return ran, mg.Fatalf(code, `running "%s %s" failed with exit code %d`, cmd, strings.Join(args, " "), code)
	}
	return ran, fmt.Errorf(`failed to run "%s %s: %v"`, cmd, strings.Join(args, " "), err)
}

func run(dir string, env map[string]string, stdout, stderr io.Writer, cmd string, args ...string) (ran bool, code int, err error) {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	for k, v := range env {
		c.Env = append(c.Env, k+"="+v)
	}
	c.Dir = dir
	c.Stderr = stderr
	c.Stdout = stdout
	c.Stdin = os.Stdin
	log.Println("exec (in '"+dir+"'):", cmd, strings.Join(args, " "))
	err = c.Run()
	return sh.CmdRan(err), sh.ExitStatus(err), err
}
