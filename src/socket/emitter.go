package socket

const (
	EMITTER_STATUS = "exec_status"
)

type Emitter interface {
	Name() string
	History() (err error, messages []Message)
	Watch() chan Message
	Close(chan Message)
}

type AbstractEmitter struct {
	pipes []chan Message
}

func NewEmitter() AbstractEmitter {
	return AbstractEmitter{
		pipes: make([]chan Message, 0),
	}
}

type Event struct {
	Room  string
	Event string
}

type Message struct {
	Data interface{}
}

func (t *AbstractEmitter) Emit(message interface{}) {
	go func() {
		for _, c := range t.pipes {
			c <- Message{
				Data: message,
			}
		}
	}()
}

func (t *AbstractEmitter) Watch() chan Message {
	c := make(chan Message)
	t.pipes = append(t.pipes, c)
	return c
}

func (t *AbstractEmitter) Close(c chan Message) {
	for i, d := range t.pipes {
		if d == c {
			t.pipes = append(t.pipes[:i], t.pipes[i+1:]...)
			return
		}
	}
}
