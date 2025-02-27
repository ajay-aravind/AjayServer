package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"syscall"
	"time"
)

type listener struct {
	port     int
	protocol string
	address  string
}

var endOfRequestLine []byte = []byte("\r\n")

func (listener listener) startWorker(tasksChannel chan Task) {

	// Define ListenConfig with platform-specific optimizations
	lc := net.ListenConfig{
		Control: controlSocketOptions, // Calls OS-specific function
	}

	// timeout for listen call to bind to the address and port;
	// if listen isn't completed in 10 seconds it will throw error
	lctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	ln, err := lc.Listen(lctx, listener.protocol, listener.address+":"+strconv.Itoa(listener.port))
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go handleControlC(ln, tasksChannel)

	log.Println("Server is running on http://localhost:" + strconv.Itoa(listener.port))

	readDuration := make(chan int64)
	var connectionId int
	// go calculateAverage(readDuration)

	for {
		connectionId++
		// Accept a connection
		conn, err := ln.Accept()

		if err != nil {
			select {
			case <-sigChan: // Check if we were asked to quit
				log.Println("Listener closed. Exiting accept loop.")
				pprof.StopCPUProfile()
				log.Println("Stopped cpu profiling")
				return // Exit main goroutine
			default:
				log.Println("Error accepting connection:", err)
				// Handle the error appropriately (e.g., log, backoff, etc.)
				// But do NOT exit the loop unless it is a fatal error
				// or if it is the error caused by closing the listener.
				if err == net.ErrClosed {
					return
				}
			}
			log.Println("Error accepting connection:", err)
		} else {
			log.Println("Accepted new connection")
		}

		if tc, ok := conn.(*net.TCPConn); ok {
			if err := tc.SetKeepAlive(true); err != nil {
				_ = tc.Close()
				log.Println("error while setting keep alive connection")
			}
		}

		// Handle the connection in a separate goroutine
		// todo create go routine pool to handle connections
		// creating go routine per request on demamd might create latency and exhaust resources
		// also we already created multiple process listening on same port, so it would be
		// fine for now to treat each process as a single threaded application.
		// But if i remove go key work from the below line i am seeing that average connection takes 2x time
		// go handleConnection(conn, readDuration)
		tasksChannel <- Task{id: connectionId, connection: &conn, readDuration: readDuration}
	}
}

func handleControlC(ln net.Listener, tasksChannel chan Task) {
	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	// Wait for termination signal
	<-sigChan
	log.Println("Got ctrl+c")
	close(tasksChannel)
	log.Println("closed tasks channel")
	err := ln.Close()
	if err != nil {
		log.Println("Server close is failed")
	}
	log.Println("closed listener")
}

func calculateAverage(readDuration chan int64) {

	var sum int64
	var count int64

	for duration := range readDuration {
		sum = (sum + duration)
		count += 1
		// log.Println("duration:" + strconv.FormatInt(duration, 10))
		// log.Println("sum:" + strconv.FormatInt(sum, 10))
		// log.Println("count:" + strconv.FormatInt(count, 10))
	}
}
