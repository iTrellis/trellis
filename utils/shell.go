package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Verbose print debug
var Verbose = false

// RunCommand executes a shell command.
func RunCommand(name string, arg ...string) error {
	if Verbose {
		cmdText := name + " " + strings.Join(arg, " ")
		fmt.Fprintln(os.Stderr, " + ", cmdText)
	}
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
