package multicast

import (
	"net"
	"testing"
)

func BenchmarkStringMapIPv4(b *testing.B) {
	ip := net.ParseIP("1.2.3.4")
	m := make(map[string]int)

	b.ResetTimer()

	var addr string
	for i := 0; i < b.N; i++ {
		addr = string(ip.To4())
		m[addr] = i
	}
}

func BenchmarkStringMapIPv6(b *testing.B) {
	ip := net.ParseIP("::FFFF:1.2.3.4")
	m := make(map[string]int)

	b.ResetTimer()

	var addr string
	for i := 0; i < b.N; i++ {
		addr = string(ip.To4())
		m[addr] = i
	}
}

func BenchmarkArrayMapIPv4(b *testing.B) {
	ip := net.ParseIP("1.2.3.4")
	m := make(map[[4]byte]int)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m[*(*[4]byte)(ip.To4())] = i
	}
}

func BenchmarkArrayMapIPv6(b *testing.B) {
	ip := net.ParseIP("::FFFF:1.2.3.4")
	m := make(map[[4]byte]int)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m[*(*[4]byte)(ip.To4())] = i
	}
}
