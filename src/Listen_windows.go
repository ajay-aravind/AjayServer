//go:build windows
// +build windows

package main

import (
	"syscall"
)

// controlSocketOptions applies OS-specific optimizations for Windows
func controlSocketOptions(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		// Windows doesn't support SO_REUSEPORT or TCP_FASTOPEN

		// This will work in windows but only for udp protocol i.e one-many
		// for protocols like tcp it won't work, even if all the proccesses can bind
		// to the same port, only one of them will recieve traffic. Sad that windows
		// won't support this. Linux is better in this way
		// err := windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1) // Enable address reuse
		// if err != nil {
		// 	fmt.Println("Error setting SO_REUSEADDR:", err)
		// }
	})
}
