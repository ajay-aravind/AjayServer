package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
)

type listener struct {
	port     int
	protocol string
	address  string
}

var endOfRequestLine []byte = []byte("\r\n")

func (listener listener) startWorker(tasksChannel chan<- Task) {

	// Define ListenConfig with platform-specific optimizations
	lc := net.ListenConfig{
		Control: controlSocketOptions, // Calls OS-specific function
	}

	// timeout for listen call to bind to the address and port;
	// if listen isn't completed in 10 seconds it will throw error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	ln, err := lc.Listen(ctx, listener.protocol, listener.address+":"+strconv.Itoa(listener.port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is running on http://localhost:" + strconv.Itoa(listener.port))
	readDuration := make(chan int64)

	go calculateAverage(readDuration)

	for {
		// Accept a connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		// todo create go routine pool to handle connections
		// creating go routine per request on demamd might create latency and exhaust resources
		// also we already created multiple process listening on same port, so it would be
		// fine for now to treat each process as a single threaded application.
		// But if i remove go key work from the below line i am seeing that average connection takes 2x time
		// go handleConnection(conn, readDuration)
		tasksChannel <- Task{connection: conn, readDuration: readDuration}
	}
}

func calculateAverage(readDuration chan int64) {

	var sum int64
	var count int64

	for duration := range readDuration {
		sum = (sum + duration)
		count += 1
		// fmt.Println("duration:" + strconv.FormatInt(duration, 10))
		// fmt.Println("sum:" + strconv.FormatInt(sum, 10))
		// fmt.Println("count:" + strconv.FormatInt(count, 10))
	}
}
