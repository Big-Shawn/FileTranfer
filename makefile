
all:server-build client-build

server-build:
	go build -C server  -o server -v main.go

client-build:
	go build -C client -o client  -v main.go
