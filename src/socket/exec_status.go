package socket

import (
	"time"

	"github.com/senfix/process-overseer/src/model"

	"github.com/senfix/process-overseer/src/service"
)

type ExecStatus interface {
	Emitter
	Emit(overseer service.Overseer)
	Run(overseer service.Overseer, delay time.Duration)
}

type execStatus struct {
	AbstractEmitter
}

func NewExecStatus() ExecStatus {
	return &execStatus{
		AbstractEmitter: NewEmitter(),
	}
}

func (t *execStatus) Name() string {
	return EMITTER_STATUS
}

func (t *execStatus) Run(overseer service.Overseer, delay time.Duration) {
	for {
		t.Emit(overseer)
		time.Sleep(time.Duration(delay))
	}
}

func (t *execStatus) Emit(overseer service.Overseer) {
	commands := overseer.GetCommands()
	data := []model.ExecStatus{}
	for _, c := range commands {
		cData := []model.ExecStatus{}
		oneActive := false
		for i := 0; i < c.Workers; i++ {
			status := overseer.Status(c.Id, i)
			state := overseer.CmdState(c.Id, i)
			intDurration := time.Duration(int64(status.Runtime * 1000 * 1000 * 1000)).Round(100 * time.Millisecond)
			cData = append(cData, model.ExecStatus{
				Exec:    c.Id,
				Worker:  i,
				State:   state.String(),
				PID:     status.PID,
				Start:   time.Unix(0, status.StartTs).Format("2006-01-02 15:04:05"),
				Stop:    time.Unix(0, status.StopTs).Format("2006-01-02 15:04:05"),
				Runtime: model.Duration{intDurration},
			})
			if state != model.INITIAL {
				oneActive = true
			}
		}

		if c.Enabled || oneActive {
			data = append(data, cData...)
		}
	}

	t.AbstractEmitter.Emit(data)

}

func (t *execStatus) History() (err error, messages []Message) {
	return nil, make([]Message, 0)
}
