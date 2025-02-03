package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type listener struct {
	port     string
	protocol string
}

func (l listener) startWorker() {

	// Define ListenConfig with platform-specific optimizations
	lc := net.ListenConfig{
		Control: controlSocketOptions, // Calls OS-specific function
	}

	// Create a context with timeout (optional)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	ln, err := lc.Listen(ctx, "tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is running on http://localhost:" + l.port)

	for {
		// Accept a connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the raw HTTP request
	reader := bufio.NewReader(conn)
	var requestLines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}

		// Stop reading when an empty line is found (end of headers)
		if line == "\r\n" {
			break
		}

		// Collect the lines of the request
		requestLines = append(requestLines, strings.TrimSpace(line))
	}

	// Print the raw request
	// fmt.Println("Raw HTTP Request:")
	// for _, line := range requestLines {
	// 	fmt.Println(line)
	// }

	// Send a simple HTTP response
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 12\r\n" +
		"\r\n" +
		"Hello World!"
	conn.Write([]byte(response))
}
