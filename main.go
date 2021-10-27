package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

const (
	maxDatagramSize = 8192
	interfaceName   = "bond0"
)

func nasdaq() func(int, []byte) {
	exp := uint64(0)
	return func(n int, buff []byte) {
		seq := binary.BigEndian.Uint64(buff[10:])
		count := binary.BigEndian.Uint16(buff[18:])

		switch exp {
		case seq:
			fmt.Println("nasdaq", seq, count)
		default:
			fmt.Println("WARN: nasdaq", exp, seq, count)
		}

		exp = seq + uint64(count)
	}
}

func nyse() func(int, []byte) {
	exp := uint32(0)
	return func(n int, buff []byte) {
		seq := binary.BigEndian.Uint32(buff[5:])
		count := buff[9]

		switch exp {
		case seq, seq + 1:
			fmt.Println("nyse", seq, count)
		default:
			fmt.Println("WARN: nyse", exp, seq, count)
		}

		exp = seq + uint32(count)
	}
}

func serveMulticastUDP(address string, handler func(int, []byte)) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		log.Fatal("SplitHostPort:", err)
	}

	conn, err := net.ListenPacket("udp4", address)
	if err != nil {
		log.Fatal("ListenPacket:", err)
	}

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Fatal("InterfaceByName:", err)
	}

	group := net.ParseIP(host)
	pconn := ipv4.NewPacketConn(conn)
	if err := pconn.JoinGroup(iface, &net.UDPAddr{IP: group}); err != nil {
		log.Fatal("JoinGroup:", err)
	}

	if err := pconn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		log.Fatal("SetContrlMessage:", err)
	}

	buff := make([]byte, maxDatagramSize)
	for {
		n, cm, _, err := pconn.ReadFrom(buff)
		if err != nil {
			log.Fatal("ReadFrom failed:", err)
		}

		if !cm.Dst.IsMulticast() {
			continue
		}

		if !cm.Dst.Equal(group) {
			continue
		}

		handler(n, buff)
	}
}

func main() {
	type addrMap struct {
		addr string
		fn   func(int, []byte)
	}

	for _, x := range []addrMap{
		{"233.46.176.8:55640", nasdaq()},
		{"224.0.89.0:40000", nyse()},
	} {
		go serveMulticastUDP(x.addr, x.fn)
	}

	<-make(chan struct{})
}
