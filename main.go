package main

import (
	"net"

	"github.com/leki75/multicast-test/log"
	"github.com/leki75/multicast-test/net/multicast"
	"github.com/leki75/multicast-test/parser"
	"go.uber.org/zap"
)

type addrHandlers struct {
	ip   net.IP
	port int
	fn   func(int) multicast.HandlerFunc
}

var iface = "en0"

func main() {
	servers := make(map[int]*multicast.UDPServer)

	var err error
	for i, addrHandler := range []addrHandlers{
		{ip: net.ParseIP("233.46.176.8"), port: 55640, fn: parser.Nasdaq},  // Nasdaq, Trade, Tape A, New York, Line A
		{ip: net.ParseIP("233.46.176.24"), port: 55640, fn: parser.Nasdaq}, // Nasdaq, Trade, Tape A, New York, Line B
		{ip: net.ParseIP("233.46.176.72"), port: 55640, fn: parser.Nasdaq}, // Nasdaq, Trade, Tape A, Chicago,  Line A
		{ip: net.ParseIP("233.46.176.88"), port: 55640, fn: parser.Nasdaq}, // Nasdaq, Trade, Tape A, Chicago,  Line B
		{ip: net.ParseIP("224.0.89.0"), port: 40000, fn: parser.Nyse},      // Nyse,   Trade, Tape A, New York, Line A
		{ip: net.ParseIP("224.0.89.128"), port: 40000, fn: parser.Nyse},    // Nyse,   Trade, Tape A, New York, Line B
	} {
		server, ok := servers[addrHandler.port]

		if !ok {
			server, err = multicast.NewUDPServer(iface, addrHandler.port)
			if err != nil {
				log.Logger.Fatal("new multicast", zap.Error(err))
			}
			servers[addrHandler.port] = server
		}

		if err := server.Listen(addrHandler.ip, addrHandler.fn(i)); err != nil {
			log.Logger.Fatal("multicast listen", zap.Error(err))
		}
	}

	var errorCh = make(chan error)
	for _, s := range servers {
		go func(s *multicast.UDPServer) {
			errorCh <- s.Serve()
		}(s)
	}

	log.Logger.Fatal("multicast listen", zap.Error(<-errorCh))
}
