package main

import (
	"client/receiver"
	"flag"
	"log"
	"net"
)

var host, path string

func init() {
	flag.StringVar(&host, "p", "127.0.0.1:8879", "server port")
	flag.StringVar(&path, "c", "./", "file storage path")
}

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}

	// receive from server
	r := receiver.New(path)
	if err := r.Read(&conn); err != nil {
		panic(err)
	}
	log.Printf("\nfile receive complete")
}
