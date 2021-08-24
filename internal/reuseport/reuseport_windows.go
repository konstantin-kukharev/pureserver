// +build windows

package reuseport

import "net"

func NewReusablePortListener(proto, addr string) (net.Listener, error) {
	return net.Listen(proto, addr)
}

func NewReusablePortPacketConn(proto, addr string) (net.PacketConn, error) {
	return net.ListenPacket(proto, addr)
}