package main

import (
	"bufio"
	"fmt"
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
	scanner := bufio.NewScanner(conn)
	// Scan() will return true until it reaches the end of a file or an error
	for scanner.Scan() {
		// print out the line
		ln := scanner.Text()
		fmt.Println(ln)
		// conn has both a reader and a writer so use it to write back to the
		// connection
		fmt.Fprintf(conn, "I heard you say: %s\n", ln)
	}
	// don't close the connection
	defer conn.Close()
}
