package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/senfix/process-overseer/src/controller/rest"
	"github.com/senfix/process-overseer/src/controller/ws"
	"github.com/senfix/process-overseer/src/server"
	"github.com/senfix/process-overseer/src/service"
	"github.com/senfix/process-overseer/src/socket"
)

var (
	ListeningPort string
	EmitterDelay  int
	StoragePath   string
)

func main() {

	flag.StringVar(&ListeningPort, "listen", "", "port to listen")
	flag.StringVar(&StoragePath, "storage", "./", "path to storage.json")
	flag.IntVar(&EmitterDelay, "emit-delay", 100, "emit status every N millisecond")
	flag.Parse()

	fmt.Printf("StoragePath: %v Listening on: %v, emitDelay: %v\n", StoragePath, ListeningPort, EmitterDelay)

	storage := service.NewStorage(StoragePath)
	overseer := service.NewOverseer(storage)

	//start overseer
	go overseer.WatchAll()
	defer overseer.StopAll()

	if len(ListeningPort) > 0 {
		//define services
		statusEmitter := socket.NewExecStatus()
		wsEmitters := []socket.Emitter{statusEmitter}
		webSocket := ws.NewSocket(wsEmitters)

		//start webserver
		s := server.NewServer(ListeningPort)
		err := s.Start(
			&rest.Commands{Storage: storage, Overseer: overseer},
			&webSocket,
		)
		if err != nil {
			panic(fmt.Sprintf("cannot server content on %v", ListeningPort))
		}
		defer s.Stop()

		//start emitting web sockets
		go statusEmitter.Run(overseer, time.Duration(EmitterDelay)*time.Millisecond)

	}

	//signal wait
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM|syscall.SIGKILL)
	//sw<-c
	switch <-sigChan {
	case syscall.SIGKILL:
		overseer.KillAll()
	case syscall.SIGTERM:
		overseer.StopAll()
	}

}
