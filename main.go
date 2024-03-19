package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
	// don't close the connection
	defer conn.Close()

	// read request
	request(conn)

	// write back to request
	respond(conn)
}

func request(conn net.Conn) {
	// counter to determine current buffer line
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// current line from buffer
		ln := scanner.Text()
		fmt.Println(ln)
		if (i == 0) {
			// split by space and return first 'word' which is the verb
			// this code seems a bit brittle really
			m := strings.Fields(ln)[0]
			fmt.Println("Method: ", m)
		}
		if (ln == "") {
			// if the buffer returns an empty string, we're done - no need to
			// iterate i
			break
		}
		i++
	}
}

func respond(conn net.Conn) {
	body := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Welcome to the desert of the real</title>
</head>
<body>
	<h1>Nobody can be told what the Matrix is...</h1>
</body>
</html>
	`

	// follow the HTTP spec in RFC 7230 specified by the Internet Engineering
	// Task Force https://datatracker.ietf.org/doc/html/rfc7230
	// return response
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}
