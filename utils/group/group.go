// Package group provides a start method for a command, concerning its group of subprocesses.
package group

import (
	"os/exec"
	"syscall"
)

// Start starts the given command in the background, and provides an
// additional kill callback for stopping the process as well as its
// subprocesses.
//
// Credits to https://stackoverflow.com/a/71714364.
//
// TODO: For Windows compatibility, refer to
// https://stackoverflow.com/a/66500411.
func Start(c *exec.Cmd, sig syscall.Signal) (kill func(), err error) {
	// Enable setpgid bit so we can kill child processes when the kill callback
	// is called
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Setpgid = true

	// Start the command
	err = c.Start()
	if err != nil {
		return func() {}, err
	}

	return func() {
		p := c.Process
		if p == nil {
			return
		}

		// Kill by negative PID to kill the whole process group, which includes
		// the top-level process we spawned as well as any subprocesses
		_ = syscall.Kill(-p.Pid, sig)
	}, nil
}
