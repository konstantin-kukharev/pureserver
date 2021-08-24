// +build !darwin,!netbsd,!freebsd,!openbsd,!dragonfly,!linux

package internal

import (
	"errors"
	"net"
	"os"
)

func (ln *listener) close() {
	if ln.ln != nil {
		ln.ln.Close()
	}
	if ln.pConn != nil {
		ln.pConn.Close()
	}
	if ln.network == "unix" {
		os.RemoveAll(ln.addr)
	}
}

func (ln *listener) system() error {
	return nil
}

func serve(events Events, listeners []*listener) error {
	return stdserve(events, listeners)
}

func reuseportListenPacket(proto, addr string) (l net.PacketConn, err error) {
	return nil, errors.New("reuseport is not available")
}

func reuseportListen(proto, addr string) (l net.Listener, err error) {
	return nil, errors.New("reuseport is not available")
}
