// +build linux darwin dragonfly freebsd netbsd openbsd

package test

import (
	"github.com/konstantin-kukharev/pureserver/internal/reuseport"
	"testing"
)

func TestNewReusablePortPacketConn(t *testing.T) {
	listenerOne, err := reuseport.NewReusablePortPacketConn("udp4", "localhost:10082")
	if err != nil {
		t.Error(err)
	}
	defer listenerOne.Close()

	listenerTwo, err := reuseport.NewReusablePortPacketConn("udp", "127.0.0.1:10082")
	if err != nil {
		t.Error(err)
	}
	defer listenerTwo.Close()

	listenerThree, err := reuseport.NewReusablePortPacketConn("udp6", "[::1]:10082")
	if err != nil {
		t.Error(err)
	}
	defer listenerThree.Close()

	listenerFour, err := reuseport.NewReusablePortListener("udp6", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerFour.Close()

	listenerFive, err := reuseport.NewReusablePortListener("udp4", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerFive.Close()

	listenerSix, err := reuseport.NewReusablePortListener("udp", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerSix.Close()
}

func TestListenPacket(t *testing.T) {
	listenerOne, err := reuseport.ListenPacket("udp4", "localhost:10082")
	if err != nil {
		t.Error(err)
	}
	defer listenerOne.Close()

	listenerTwo, err := reuseport.ListenPacket("udp", "127.0.0.1:10082")
	if err != nil {
		t.Error(err)
	}
	defer listenerTwo.Close()

	listenerThree, err := reuseport.ListenPacket("udp6", "[::1]:10082")
	if err != nil {
		t.Error(err)
	}
	defer listenerThree.Close()

	listenerFour, err := reuseport.ListenPacket("udp6", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerFour.Close()

	listenerFive, err := reuseport.ListenPacket("udp4", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerFive.Close()

	listenerSix, err := reuseport.ListenPacket("udp", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerSix.Close()
}

func BenchmarkNewReusableUDPPortListener(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listener, err := reuseport.NewReusablePortPacketConn("udp4", "localhost:10082")

		if err != nil {
			b.Error(err)
		} else {
			listener.Close()
		}
	}
}
