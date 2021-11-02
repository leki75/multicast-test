package main

import (
	"flag"
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

var iface string

func init() {
	flag.StringVar(&iface, "iface", "bond0", "bind multicast addresses to this interface")
}

func main() {
	flag.Parse()

	servers := make(map[int]*multicast.UDPServer)

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
			var err error
			server, err = multicast.NewUDPListener(iface, addrHandler.port)
			if err != nil {
				log.Logger.Fatal("new multicast", zap.Error(err))
			}
			servers[addrHandler.port] = server
		}

		if err := server.Join(addrHandler.ip, addrHandler.fn(i)); err != nil {
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
