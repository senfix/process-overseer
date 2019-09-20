package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/senfix/process-overseer/src/socket"
)

func NewSocket(emitters []socket.Emitter) WebSocket {
	return WebSocket{
		hubs:     map[string]*Hub{},
		emitters: emitters,
	}

}

type WebSocket struct {
	hub      *Hub
	hubs     map[string]*Hub
	emitters []socket.Emitter
}

func (t *WebSocket) EmitterHandler(emitter socket.Emitter, hub *Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()

		//handle history
		err, history := emitter.History()
		if err != nil {
			fmt.Printf("[%v] history: %v\n", emitter.Name(), err)
		}

		for _, d := range history {
			data, err := json.Marshal(d.Data)
			if err != nil {
				continue
			}
			client.send <- data
		}

	}
}

func (t *WebSocket) Register(root *mux.Router) {
	p := root.PathPrefix("/ws").Subrouter()

	for _, emitter := range t.emitters {
		hub := newHub()
		go hub.run()

		data := emitter.Watch()
		go func() {
			for {
				m := <-data
				msg, err := json.Marshal(m.Data)

				if err != nil {
					continue
				}
				hub.broadcast <- msg
			}

		}()
		p.HandleFunc(fmt.Sprintf("/%v", emitter.Name()), t.EmitterHandler(emitter, hub)).Methods("GET")
	}

	p.StrictSlash(true)
}
