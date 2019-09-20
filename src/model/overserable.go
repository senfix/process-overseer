package model

type Overserable struct {
	Command Command
	Runners map[int]Runner
}

func (t *Overserable) Init() Runner {
	return Runner{
		command: &t.Command,
		exec: NewCmd(t.Command.Exec, t.Command.Args, Options{
			Dir:       t.Command.WorkDir,
			Buffered:  false,
			Streaming: true,
		}),
	}
}

func (t *Overserable) GetRunner(worker int) Runner {
	runner, _ := t.Runners[worker]
	return runner
}

func (t *Overserable) SaveRunner(worker int, runner Runner) {
	t.Runners[worker] = runner
}
