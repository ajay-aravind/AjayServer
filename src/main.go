package main

import "runtime"

const PortToBind = 8080

// this is bad practice to bind on all zeros address. Essentially this listens on
// all ips of network interfaces avialble on the machine. Ideally we would want to
// listen on particular network interface private/public depending on the use case
const AddressToBind = "0.0.0.0"

// number of processes listening on the port
const ProcessCount = 10

// go Routine pool count
var WorkerPoolCount int = runtime.NumCPU()

const Protocol = "tcp"

func main() {

	// single thread vs multi thread
	// when using single process i see that 1000 request with 50 clients taking average of 9 seconds
	// with 5 worker threads, i see that it is taking 5 seconds
	singleThreadedMain(PortToBind, AddressToBind, Protocol)
	// singleThreadedMain(PortToBind, ProcessCount, AddressToBind, Protocol)
}

func singleThreadedMain(PortToBind int, AddressToBind string, Protocol string) {
	tasksChannel := InitGoRoutinePool(WorkerPoolCount)
	var server listener = listener{port: PortToBind, protocol: Protocol, address: AddressToBind}
	server.startWorker(tasksChannel)
}
