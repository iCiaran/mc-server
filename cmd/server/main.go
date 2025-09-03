package main

import (
	"log"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("tcp4", ":25565")
	if err != nil {
		log.Printf("Error listening on tcp: %v\n", err)
		os.Exit(1)
	}

	log.Println("Listening on tcp:", ln.Addr())

	for {
		log.Println("Accepting connection")
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}

		go func() {
			err := conn.Close()
			if err != nil {
				log.Printf("Error closing connection: %v\n", err)
			}
		}()
	}
}
