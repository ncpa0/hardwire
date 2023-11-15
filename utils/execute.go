package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

type ExecuteOptions struct {
	Wd  string
	Env map[string]string
}

type ExecuteResult struct {
	Stdout string
	Stderr string
	Err    error
}

func Execute(cmd string, args []string, options *ExecuteOptions) *ExecuteResult {
	command := exec.Command(
		cmd,
		args...,
	)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if options != nil {
		if options.Wd != "" {
			command.Dir = options.Wd
		}
		if options.Env != nil {
			command.Env = command.Environ()
			for key, value := range options.Env {
				command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	err := command.Run()

	return &ExecuteResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
}
