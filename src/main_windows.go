//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func multiThreadedMain(PortToBind int, ProcessCount int, AddressToBind string, Protocol string) {
	var workerCount int = ProcessCount
	// Check if this is a child process (forked)
	if len(os.Args) > 1 && os.Args[1] == "worker" {
		tasksChannel := InitGoRoutinePool(WorkerPoolCount)
		startWorkerProcess(PortToBind, AddressToBind, Protocol, tasksChannel)
		return
	}

	// Parent process: Spawning multiple listener processes
	for i := 0; i < workerCount; i++ {
		cmd := exec.Command(os.Args[0], "worker") // Fork a new process
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			log.Fatalf("Failed to start worker %d: %v", i, err)
		}
		fmt.Printf("Started worker %d (PID: %d)\n", i, cmd.Process.Pid)
	}

}

func startWorkerProcess(PortToBind int, AddressToBind string, Protocol string, tasksChannel chan<- Task) {
	var server listener = listener{port: PortToBind, protocol: Protocol, address: AddressToBind}
	server.startWorker(tasksChannel)
}
