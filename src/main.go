package main

import (
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
)

const PortToBind = 8080

// this is bad practice to bind on all zeros address. Essentially this listens on
// all ips of network interfaces avialble on the machine. Ideally we would want to
// listen on particular network interface private/public depending on the use case
const AddressToBind = "0.0.0.0"

// number of processes listening on the port
const ProcessCount = 10

// go Routine pool count
// since we are testing with long lived connections it doesn't make sense to handle multiple connections in single go routine
// we can try to handle multiple go http requests in fixed pool of go routines though. This is still something to do
var WorkerPoolCount int = 600

const Protocol = "tcp"

const BufferPoolSize int = 20 * 1024 //4 x 200 x 1024 x 1024 B =  800MB

func main() {
	// Start profiling
	f, err := os.Create("myprogram.prof")
	if err != nil {
		log.Println("cloud not create profile")
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Println("Couldn't start profiling")
	}
	// single thread vs multi thread
	// when using single process i see that 1000 request with 50 clients taking average of 9 seconds
	// with 5 worker threads, i see that it is taking 5 seconds
	singleThreadedMain(PortToBind, AddressToBind, Protocol)
	// singleThreadedMain(PortToBind, ProcessCount, AddressToBind, Protocol)
}

func singleThreadedMain(PortToBind int, AddressToBind string, Protocol string) {
	tasksChannel := InitGoRoutinePool(WorkerPoolCount)
	InitBufferPools(BufferPoolSize)
	var server listener = listener{port: PortToBind, protocol: Protocol, address: AddressToBind}
	server.startWorker(tasksChannel)
}
