package model

import (
	"syscall"
)

type Runner struct {
	command *Command
	exec    *Cmd
	status  <-chan Status
}

func (t *Runner) Start() {
	if t.exec.IsFinalState() {
		t.exec = t.exec.Clone()
	}
	t.status = t.exec.Start()
}

func (t *Runner) Status() Status {
	return t.exec.Status()
}

func (t *Runner) StatusChan() <-chan Status {
	return t.status
}

func (t *Runner) State() CmdState {
	return t.exec.State
}
func (t *Runner) IsFinalState() bool {
	return t.exec.IsFinalState()
}
func (t *Runner) IsInitialState() bool {
	return t.exec.IsInitialState()
}

func (t *Runner) Stdout() chan string {
	return t.exec.Stdout
}

func (t *Runner) Stderr() chan string {
	return t.exec.Stderr
}

func (t *Runner) Close() (err error) {
	return t.exec.Stop()
}

func (t *Runner) Kill() {
	t.exec.Signal(syscall.SIGKILL)
}
