all: build run

build:
	go build -o bin/vss src/main.go

run:
	bin/vss -router -runner -log-enable -log-path="logs/logs.txt" -log-level=debug
