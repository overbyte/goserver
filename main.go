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
}

func request(conn net.Conn) {
	// counter to determine current buffer line
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// current line from buffer
		ln := scanner.Text()
		fmt.Println(ln)
		// print out verb, probably GET right now
		if (i == 0) {
			// split by space and return first 'word' which is the verb
			// this code seems a bit brittle really but it's good to learn the
			// why before we dive into a concrete solution
			mux(conn, ln)
		}
		if (ln == "") {
			// if the buffer returns an empty string, we're done - no need to
			// iterate i
			break
		}
		i++
	}
}

func mux(conn net.Conn, ln string) {
	f := strings.Fields(ln)
	method := f[0]
	uri := f[1]
	fmt.Println("Method:", method)
	fmt.Println("URI:", uri)

	// we only need to respond to GET requests right now
	if (method == "GET") {
		// setup routes
		switch uri {
			case "/about":
				aboutPage(conn)
				break
			case "/allandt":
				mePage(conn)
				break
			default:
				fmt.Println("No route found", uri)
		}
	}
}

func mePage(conn net.Conn) {
	body := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Welcome to the desert of the real</title>
</head>
<body>
	<h1>Allandt Page</h1>
</body>
</html>
	`

	respond(conn, body)
}

func aboutPage(conn net.Conn) {
	body := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Welcome to the desert of the real</title>
</head>
<body>
	<h1>About Page</h1>
</body>
</html>
	`
	respond(conn, body)
}


func respond(conn net.Conn, body string) {
	// follow the HTTP spec in RFC 7230 specified by the Internet Engineering
	// Task Force https://datatracker.ietf.org/doc/html/rfc7230
	// return response
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}
