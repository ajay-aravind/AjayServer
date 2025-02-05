package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

type Task struct {
	id           int
	connection   net.Conn
	readDuration chan int64
}

// Worker function that processes tasks
func worker(id int, tasks <-chan Task) {
	for task := range tasks {
		handleConnection(task.connection, task.readDuration)
	}
}

func createWorkerPool(workerCount int, tasks <-chan Task) {
	// Start worker goroutines
	for i := 1; i <= workerCount; i++ {
		go worker(i, tasks)
	}
}

func InitGoRoutinePool(poolCount int) chan Task {
	tasks := make(chan Task, poolCount)
	createWorkerPool(poolCount, tasks)
	return tasks
}

func handleConnection(conn net.Conn, readDuration chan int64) {
	defer conn.Close()

	//todo understand what is this setting and why this is getting used,
	// prior to using this i am getting EOF errors, when i try to benchmark my server with wrk aginst fasthttp
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Read the raw HTTP request
	reader := bufio.NewReader(conn)

	// out the below two, ideally readBytes should give less latency
	// but for some weird reason i am seeing both performe almost equally
	duration := readBytes(reader)
	readDuration <- duration.Microseconds()
	// readString(reader)

	// Send a simple HTTP response
	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!"
	conn.Write([]byte(response))

}

func readString(reader *bufio.Reader) {

	var requestLines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}

		if line == "\r\n" {
			break
		}
		// Collect the lines of the request
		requestLines = append(requestLines, strings.TrimSpace(line))
	}
	// Print the raw request
	// for _, line := range requestLines {
	// 	fmt.Println(line)
	// }
}

func readBytes(reader *bufio.Reader) time.Duration {
	start := time.Now() // Start the stopwatch

	// output buffer
	var buf bytes.Buffer

	for {
		//todo does this needs to be "\n" or "\r\n" or something like that, needs to decide
		bytesOfLine, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading request:", err)
			return time.Nanosecond
		}
		// Stop reading when an empty line is found (end of headers)
		if bytes.Equal(bytesOfLine, endOfRequestLine) {
			break
		}

		// Collect the lines of the request
		buf.Write(bytesOfLine)
	}
	// Print the raw request
	// fmt.Println("Raw HTTP Request:")
	// fmt.Println(buf.String())
	elapsed := time.Since(start)
	return elapsed
}
