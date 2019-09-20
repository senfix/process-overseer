package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/senfix/process-overseer/src/model"
)

const (
	storage_file = "storage.json"
)

type Storage interface {
	Persist()
	GetCommand(id string) (err error, commands model.Command)
	GetCommands() (commands []model.Command)
	SaveCommands(commands []model.Command)
	SaveCommand(command model.Command)
}

func NewStorage(path string) Storage {
	s := &storage{path: path}
	s.loadStorage()
	return s
}

type storage struct {
	storage model.Storage
	path    string
}

func (t *storage) loadStorage() {
	jsonFile, err := os.Open(path.Join(t.path, storage_file))
	// if we os.Open returns an error then handle it
	if err != nil {
		//panic(err)
	}
	defer jsonFile.Close()

	storage := model.Storage{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &storage)
	if err != nil {
		//panic(err)
	}
	t.storage = storage
}

func (t *storage) Persist() {
	file, _ := json.MarshalIndent(t.storage, "", " ")
	_ = ioutil.WriteFile(path.Join(t.path, storage_file), file, 0644)
}

func (t *storage) GetCommand(id string) (err error, command model.Command) {
	commands := t.GetCommands()
	for _, command := range commands {
		if command.Id == id {
			return nil, command
		}
	}
	err = errors.New("cannot find command by id")
	return
}

func (t *storage) GetCommands() (commands []model.Command) {
	return t.storage.Commands
}

func (t *storage) SaveCommands(commands []model.Command) {
	t.storage.Commands = commands
}

func (t *storage) SaveCommand(command model.Command) {
	commands := t.GetCommands()
	saved := false
	for idx, c := range commands {
		if c.Id == command.Id {
			commands[idx] = command
			saved = true
			break
		}
	}

	if saved == false {
		commands = append(commands, command)
	}

	t.SaveCommands(commands)

	return
}
