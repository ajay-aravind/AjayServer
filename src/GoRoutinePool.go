package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
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
		// logrus.Debug("starting to handle task:" + strconv.Itoa(task.id))
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
			logrus.Debug("Issue while writing the response" + err.Error())
		} else {
			// logrus.Debug("wrote response bytes:" + strconv.Itoa(n) + "connection id:" + strconv.Itoa(conectionId))
		}

		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		readerPool.Put(reader)
	}
}

func readBytes(reader *bufio.Reader) error {
	// output buffer
	var buf bytes.Buffer
	var request HttpRequest

	for {
		//todo does this needs to be "\n" or "\r\n" or something like that, needs to decide
		bytesOfLine, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				logrus.Debug("client closed connection: EOF")
			} else if os.IsTimeout(err) {
				logrus.Debug("Read timeout exceeded, closing connection")
			} else {
				logrus.Debug("Error reading request:", err)
			}
			return err
		} else {
			// logrus.Debug("read one line")
		}
		// Stop reading when an empty line is found (end of headers)
		if bytes.Equal(bytesOfLine, endOfRequestLine) {
			logrus.Debug("Done reading request headers")
			logrus.Debug("HTTP Request till now:")
			logrus.Debug(buf.String())

			parser := HttpRequestParserString{}
			request = parser.ParseTillRequestHeaders(buf.Bytes())
			break
		}

		// Collect the lines of the request
		buf.Write(bytesOfLine)
		logrus.Trace("bytes so far: ", string(bytesOfLine))
	}
	// Print the raw request

	err := request.ValidateRequestTillHeaders()

	if err != nil {
		logrus.Error("Invalid request")
		return nil
	}

	if request.IsRequestAllowedToHaveBody() {
		logrus.Debug("Request may have content")
		contentLength, err := strconv.Atoi(request.Headers["content-length"])
		if err != nil {
			log.Fatal("error in formatting content-lenght")
		}

		var body []byte
		if contentLength > 0 {
			body = readRequestBody(reader, contentLength)
		}
		parser := HttpRequestParserString{}
		parser.parseHttpRequestBody(body, &request)
		request.PrintContent()
	}

	// call the custom handler provided by developer

	return nil
}

func readRequestBody(reader *bufio.Reader, contentLength int) []byte {

	logrus.Debug("started reading request content")
	contentBuf := make([]byte, contentLength)
	_, err := io.ReadFull(reader, contentBuf)
	if err != nil {
		logrus.Warn("Error while reading content: ", err)
	}
	logrus.Debug("Read request content bytes: ", contentLength)
	logrus.Debug("RequestContent: ", string(contentBuf))
	return contentBuf
}
