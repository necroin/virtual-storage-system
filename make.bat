go build -o bin\vss.exe src\main.go
bin\vss.exe -router -runner -log-enable -log-path="logs/logs.txt" -log-level=debug -token=abc %*