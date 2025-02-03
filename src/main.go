package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, World!")
	// single thread vs multi thread
	// when using single process i see that 1000 request with 50 clients taking average of 9 seconds
	// with 5 worker threads, i see that it is taking 5 seconds
	multiThreadedMain()
	// singleThreadedMain()
	time.Sleep(200 * time.Second)
}

func singleThreadedMain() {
	var server listener = listener{port: ":8080", protocol: "tcp"}
	server.startWorker()
}
