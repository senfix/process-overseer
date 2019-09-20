package server

import "github.com/gorilla/mux"

type Controller interface {
	Register(root *mux.Router)
}
