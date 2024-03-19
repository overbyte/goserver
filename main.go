package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	// start listening on port 8080
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}

	// don't immediately close the connection
	defer li.Close()

	// loop - respond to connection request by accepting it and passing it to a
	// handler function
	for {
		// on a connection, accept
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}

		// use goroutine to respond
		go handle(conn)
	}

}

// the handler function accepts a connection which is both a Reader and a Writer
func handle(conn net.Conn) {
	io.WriteString(conn, "\nHello. Welcome to the Matrix\n")
	fmt.Fprintln(conn, "This is the desert of the real")
	fmt.Fprintf(conn, "%v", "this is percent v")
	conn.Close()
}
