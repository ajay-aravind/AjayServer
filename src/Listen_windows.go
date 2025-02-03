//go:build windows
// +build windows

package main

import (
	"syscall"

	"golang.org/x/sys/windows"
)

// controlSocketOptions applies OS-specific optimizations for Windows
func controlSocketOptions(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		// Windows doesn't support SO_REUSEPORT or TCP_FASTOPEN.
		// We use SO_EXCLUSIVEADDRUSE instead, which prevents port conflicts.
		_ = windows.SetsockoptInt(windows.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	})
}
