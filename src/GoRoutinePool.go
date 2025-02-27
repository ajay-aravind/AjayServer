package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Task struct {
	id           int
	connection   *net.Conn
	readDuration chan int64
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
		// log.Println("starting to handle task:" + strconv.Itoa(task.id))
		handleConnection(task.connection, task.readDuration, task.id)
	}
}

func createWorkerPool(workerCount int, tasks <-chan Task) {
	// Start worker goroutines
	for i := 1; i <= workerCount; i++ {
		go worker(i, tasks)
	}
}

func InitBufferPools(poolSize int) {
	for i := 0; i < poolSize; i++ {
		// Create a new buffer and put it into the pool.
		readerPool.Put(bufio.NewReaderSize(nil, 4096))
	}
	log.Printf("Preallocated %d buffers of %d bytes each\n", poolSize, 4096)
}

func InitGoRoutinePool(poolCount int) chan Task {
	tasks := make(chan Task, poolCount)
	createWorkerPool(poolCount, tasks)
	return tasks
}

func handleConnection(connPointer *net.Conn, readDuration chan int64, conectionId int) {

	conn := *connPointer

	// Ensure conn is a *net.TCPConn
	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		// Disable Nagle's Algorithm (send packets immediately)
		tcpConn.SetNoDelay(true)
		// log.Printf("Disabled nagles algorithm on the connection")
	}

	for {
		// Get a reused bufio.Reader from the pool
		reader := readerPool.Get().(*bufio.Reader)
		reader.Reset(conn) // Reset reader for new connection

		// out the below two, ideally readBytes should give less latency
		// but for some weird reason i am seeing both performe almost equally
		err := readBytes(reader)
		if err != nil {
			conn.Close()
			return
		}

		// Send a simple HTTP response
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!"
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Println("Issue while writing the response" + err.Error())
		} else {
			// log.Println("wrote response bytes:" + strconv.Itoa(n) + "connection id:" + strconv.Itoa(conectionId))
		}

		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		readerPool.Put(reader)
	}
}

func readBytes(reader *bufio.Reader) error {
	// output buffer
	var buf bytes.Buffer

	for {
		//todo does this needs to be "\n" or "\r\n" or something like that, needs to decide
		bytesOfLine, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("client closed connection: EOF")
			} else if os.IsTimeout(err) {
				log.Println("Read timeout exceeded, closing connection")
			} else {
				log.Println("Error reading request:", err)
			}
			return err
		} else {
			// log.Println("read one line")
		}
		// Stop reading when an empty line is found (end of headers)
		if bytes.Equal(bytesOfLine, endOfRequestLine) {
			log.Println("request read end")
			request := HttpRequest{}
			request.ParseHttpRequest(buf.Bytes())
			break
		}

		// Collect the lines of the request
		buf.Write(bytesOfLine)
	}
	// Print the raw request
	// log.Println("Raw HTTP Request:")
	// log.Println(buf.String())
	return nil
}

func readString(reader *bufio.Reader) {

	var requestLines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading request:", err)
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
	// 	log.Println(line)
	// }
}
