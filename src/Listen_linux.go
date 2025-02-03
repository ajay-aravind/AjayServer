//go:build !windows
// +build !windows

package main

import (
	"syscall"
)

// controlSocketOptions applies OS-specific optimizations for Linux/macOS
func controlSocketOptions(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		// Enable SO_REUSEADDR and SO_REUSEPORT (allow port reuse)
		_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)

		// Enable TCP_FASTOPEN to reduce latency in new connections (Linux only)
		_ = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_FASTOPEN, 1)
	})
}
