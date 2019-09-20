.DEFAULT_GOAL := help
OVERRIDING=default

help:

build:
	go build -o bin/cmd/process_overseer cmd/process-overseer/main.go
	chmod 766 bin/cmd/process_overseer
