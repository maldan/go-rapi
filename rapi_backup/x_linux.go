package rapi_backup

import (
	"os/exec"
	"strings"
)

func Exec(args ...string) (string, string, error) {
	c := exec.Command(args[0], args[1:]...)
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	c.Stdout = stdout
	c.Stderr = stderr

	err := c.Run()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	err = c.Process.Release()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}
