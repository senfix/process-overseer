package service

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/senfix/process-overseer/src/model"
)

type Overseer interface {
	WatchAll()
	StopAll()
	KillAll()
	GetCommands() (command []model.Command)
	RegisterNew(watch model.Command)
	UpdateCommand(watch model.Command)
	Status(command string, worker int) model.Status
	CmdState(command string, worker int) model.CmdState
}

type overseer struct {
	*sync.Mutex
	overserable map[string]model.Overserable
}

func NewOverseer(storage Storage) Overseer {

	o := &overseer{
		Mutex:       &sync.Mutex{},
		overserable: map[string]model.Overserable{},
	}

	for _, command := range storage.GetCommands() {
		o.RegisterNew(command)
	}

	return o
}

func (t *overseer) WatchAll() {
	for {
		t.Lock()
		for j, o := range t.overserable {
			for i, runner := range o.Runners {
				if !o.Command.Enabled {
					break
				}
				if runner.IsInitialState() {
					fmt.Printf("start %v: %v\n", o.Command.Id, i)
					runner.Start()
				}
				if runner.IsFinalState() && o.Command.KeepAlive {
					fmt.Printf("retry %v: %v in %v\n", o.Command.Id, i, o.Command.RetryDelay)
					time.Sleep(o.Command.RetryDelay.Duration)
					fmt.Printf("start %v: %v\n", o.Command.Id, i)
					runner.Start()
				}
				t.overserable[j].Runners[i] = runner
			}
		}
		t.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func (t *overseer) stopCommand(overserable model.Overserable) {
	for id, runner := range overserable.Runners {
		err := runner.Close()
		if err != nil {
			fmt.Printf("stop %v: %v\n", overserable.Command.Id, id)
			runner.Kill()
		}
	}

}

func (t *overseer) StopAll() {
	t.Lock()
	defer t.Unlock()
	for _, o := range t.overserable {
		t.stopCommand(o)
	}

}

func (t *overseer) KillAll() {
	t.Lock()
	defer t.Unlock()
	for _, o := range t.overserable {
		for _, runner := range o.Runners {
			runner.Kill()
		}
	}

}

func (t *overseer) runCommand(c model.Command, id int, runner model.Runner) {
	for {
		select {
		case line := <-runner.Stdout():
			fmt.Sprintln(line)
		case line := <-runner.Stderr():
			fmt.Println(line)
			_, err := fmt.Fprintln(os.Stderr, line)
			if err != nil {
				panic(err)
			}
		case <-runner.StatusChan():
		}
	}
}

func (t *overseer) GetCommands() (command []model.Command) {
	t.Lock()
	defer t.Unlock()
	commands := make([]model.Command, len(t.overserable))
	i := 0
	for _, o := range t.overserable {
		commands[i] = o.Command
		i++
	}
	return commands
}

func (t *overseer) RegisterNew(command model.Command) {
	t.Lock()
	defer t.Unlock()
	overserable := model.Overserable{
		Command: command,
		Runners: map[int]model.Runner{},
	}

	runners := map[int]model.Runner{}
	for i := 0; i < command.Workers; i++ {
		runner := overserable.Init()
		go t.runCommand(command, i, runner)
		runners[i] = runner
	}

	overserable.Runners = runners
	t.overserable[command.Id] = overserable
}

func (t *overseer) UpdateCommand(command model.Command) {
	t.Lock()
	defer t.Unlock()
	old, ok := t.overserable[command.Id]
	//does not exists yet
	if ok == false {
		t.RegisterNew(command)
		return
	}

	if old.Command.Exec != command.Exec ||
		old.Command.WorkDir != command.WorkDir ||
		reflect.DeepEqual(old.Command.Args, command.Args) {
		//reinit
	}

	//turn switch
	if old.Command.Enabled != command.Enabled {
		if command.Enabled {
			old.Command = command
		} else {
			t.stopCommand(old)
		}
	}

	for old.Command.Workers != command.Workers {
		//spawn more
		if old.Command.Workers < command.Workers {
			runner := old.Init()
			idx := len(old.Runners)
			old.Runners[idx] = runner
			go t.runCommand(command, idx, runner)
		}
		//kill some
		if old.Command.Workers > command.Workers {
			idx := len(old.Runners) - 1
			runner := old.Runners[idx]
			runner.Kill()
			delete(old.Runners, idx)
		}
		old.Command.Workers = len(old.Runners)
	}

	old.Command = command
	t.overserable[command.Id] = old
}

func (t *overseer) Status(command string, worker int) model.Status {
	t.Lock()
	defer t.Unlock()
	overserable, ok := t.overserable[command]
	if ok == false {
		return model.Status{}
	}

	runner, ok := overserable.Runners[worker]
	if ok == false {
		return model.Status{}
	}

	return runner.Status()

}

func (t *overseer) CmdState(command string, worker int) model.CmdState {
	t.Lock()
	defer t.Unlock()
	overserable, ok := t.overserable[command]
	if ok == false {
		return model.INITIAL
	}

	runner, ok := overserable.Runners[worker]
	if ok == false {
		return model.INITIAL
	}

	return runner.State()
}
