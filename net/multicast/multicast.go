package multicast

import (
	"errors"
	"fmt"
	"net"
	"runtime"

	"github.com/leki75/multicast-test/net/reuse"
	"golang.org/x/net/ipv4"
)

const (
	maxDatagramSize = 8192
	soRecvBufSize   = 16_777_216
)

type UDPServer struct {
	conn         *ipv4.PacketConn
	iface        *net.Interface
	addrHandlers map[[4]byte]HandlerFunc
}

type HandlerFunc func(int, []byte)

func NewUDPServer(ifname string, port int) (*UDPServer, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}

	if iface.Flags&net.FlagMulticast == 0 {
		return nil, errors.New("interface does not support multicast")
	}

	listen, err := reuse.ListenPacket("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	if runtime.GOOS == "linux" {
		// Sets SO_RCVBUF to a high value to eliminate kernel RcvbufErrors in
		// Udp section of /proc/net/snmp
		if err := listen.(*net.UDPConn).SetReadBuffer(soRecvBufSize); err != nil {
			return nil, err
		}
	}

	conn := ipv4.NewPacketConn(listen)
	if err := conn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		return nil, err
	}

	return &UDPServer{
		iface:        iface,
		conn:         conn,
		addrHandlers: make(map[[4]byte]HandlerFunc),
	}, nil
}

func (us *UDPServer) Listen(ip net.IP, handler HandlerFunc) error {
	if ip.To4() == nil {
		return errors.New("non IPv4 address")
	}

	if err := us.conn.JoinGroup(us.iface, &net.UDPAddr{IP: ip}); err != nil {
		return err
	}

	var addr [4]byte
	copy(addr[:], ip.To4())
	us.addrHandlers[addr] = handler

	return nil
}

func (us *UDPServer) Serve() error {
	buff := make([]byte, maxDatagramSize)

	var ip net.IP
	for {
		n, cm, _, err := us.conn.ReadFrom(buff)
		if err != nil {
			return err
		}

		if !cm.Dst.IsMulticast() {
			continue
		}

		// Only handle IPv4 and IPv4-mapped IPv6 addresses
		ip = cm.Dst.To4()
		if ip == nil {
			continue
		}

		// Supported from Go 1.17
		// https://tip.golang.org/ref/spec#Conversions_from_slice_to_array_pointer
		handler, ok := us.addrHandlers[*(*[4]byte)(ip)]
		if !ok {
			continue
		}

		handler(n, buff)
	}
}
