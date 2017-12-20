package conn

import (
	"fmt"
	"net"
)

var (
	_ fmt.Stringer
)

type tcpListenConn struct {
	net.Listener
	filter func(string) bool
}

func (tl *tcpListenConn) Accept() (net.Conn, error) {
	client, err := tl.Listener.Accept()
	if err != nil {
		return client, err
	}
	rAddr := client.RemoteAddr()
	//fmt.Printf("client remote addr: <%v>\n", rAddr.String())

	// call filter, but ignore the result para este momento
	if !tl.filter(rAddr.String()) {
		client.Close()
		return client, nil
	}
	return client, nil
}

func NewTcpListener(addr string, filter func(string) bool) (*tcpListenConn, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &tcpListenConn{
		Listener: ln,
		filter:   filter,
	}, nil
}
