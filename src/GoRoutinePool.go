package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Task struct {
	id           int
	connection   net.Conn
	readDuration chan int64
}

// Pool for reusing buffers
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer) // Preallocate buffer
	},
}

// Pool for reusing bufio.Reader and its buffer
var readerPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReaderSize(nil, 4096) // Preallocate 4KB buffer
	},
}

// Worker function that processes tasks
func worker(id int, tasks <-chan Task) {
	for task := range tasks {
		handleConnection(task.connection, task.readDuration, id)
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

func handleConnection(conn net.Conn, readDuration chan int64, goRoutineId int) {
	defer conn.Close()
	//todo understand what is this setting and why this is getting used,
	// prior to using this i am getting EOF errors, when i try to benchmark my server with wrk aginst fasthttp
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Get a reused bufio.Reader from the pool
	reader := readerPool.Get().(*bufio.Reader)
	reader.Reset(conn)           // Reset reader for new connection
	defer readerPool.Put(reader) // Return reader to pool

	// Get a reused buffer from the pool
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()               // Clear the buffer
	defer bufferPool.Put(buf) // Return buffer to pool

	// out the below two, ideally readBytes should give less latency
	// but for some weird reason i am seeing both performe almost equally
	duration := readBytes(reader)
	readDuration <- duration.Microseconds()
	// readString(reader)

	// Send a simple HTTP response
	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!"
	conn.Write([]byte(response))
}

func readBytes(reader *bufio.Reader) time.Duration {
	start := time.Now() // Start the stopwatch

	// output buffer
	var buf bytes.Buffer

	for {
		//todo does this needs to be "\n" or "\r\n" or something like that, needs to decide
		bytesOfLine, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("client closed connection: EOF")
			} else {
				fmt.Println("Error reading request:", err)
			}
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
