package rest

import (
	"encoding/json"
	"net/http"

	"github.com/senfix/process-overseer/src/server"

	"github.com/gorilla/mux"

	"github.com/senfix/process-overseer/src/model"

	"github.com/senfix/process-overseer/src/service"
)

type Commands struct {
	Storage  service.Storage
	Overseer service.Overseer
}

func (t *Commands) Register(root *mux.Router) {
	root.HandleFunc("/rest/command", t.GetCommands).Methods("GET")
	root.HandleFunc("/rest/command", t.SaveCommand).Methods("POST")
	root.HandleFunc("/rest/command/{id}", t.GetCommand).Methods("GET")
}

func (t *Commands) GetCommands(w http.ResponseWriter, req *http.Request) {
	commands := t.Storage.GetCommands()

	w.Header().Add("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(commands)
	if err != nil {
		server.EmitError(w, http.StatusUnprocessableEntity, err)
	}

}

func (t *Commands) GetCommand(w http.ResponseWriter, req *http.Request) {
	id, err := server.GetParamString(req, "id")
	if err != nil {
		server.EmitError(w, http.StatusUnprocessableEntity, err)
		return
	}

	err, cmd := t.Storage.GetCommand(id)
	if err != nil {
		server.EmitError(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Add("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(cmd)
	if err != nil {
		server.EmitError(w, http.StatusUnprocessableEntity, err)
	}

}

func (t *Commands) SaveCommand(w http.ResponseWriter, req *http.Request) {
	command := model.Command{}
	err := server.Decode(w, req.Body, &command)
	if err != nil {
		server.EmitError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if command.Workers < 0 {
		command.Workers = 0
	}

	t.Storage.SaveCommand(command)
	t.Storage.Persist()
	t.Overseer.UpdateCommand(command)

}
