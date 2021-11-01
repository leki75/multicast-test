package reuse

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func control(network, address string, c syscall.RawConn) error {
	var err error

	c.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
		if err != nil {
			return
		}

		err = syscall.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if err != nil {
			return
		}
	})

	return err
}

func ListenPacket(network, address string) (net.PacketConn, error) {
	config := net.ListenConfig{Control: control}
	return config.ListenPacket(context.Background(), network, address)
}
