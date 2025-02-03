//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func multiThreadedMain() {
	var workerCount int = 10
	// Check if this is a child process (forked)
	if len(os.Args) > 1 && os.Args[1] == "worker" {
		startWorkerProcess()
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

func startWorkerProcess() {
	var server listener = listener{port: ":8080", protocol: "tcp"}
	server.startWorker()
}
