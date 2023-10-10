all: build run

build:
	go build -o bin/vss src/main.go

run:
	bin/vss
