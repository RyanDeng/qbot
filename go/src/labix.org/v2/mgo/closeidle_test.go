package mgo

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"labix.org/v2/mgo/bson"
)

type M bson.M

type logger struct {
}

func (l *logger) Output(calldepth int, s string) error {

	fmt.Printf("LOG:%v => %v\n", calldepth, s)
	return nil
}

func startTcpServer() string {

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		l, err = net.Listen("tcp6", "[::1]:0")
		if err != nil {
			panic("tcpServer: listen local failed: " + err.Error())
		}
	}
	go func() {
		defer l.Close()
		for {
			rw, e := l.Accept()
			if e != nil {
				fmt.Printf("tcpServer: Accept failed %v\n", e)
				return
			}
			go func() {
				b := make([]byte, 4096)
				for {
					_, e := rw.Read(b)
					if e != nil {
						fmt.Printf("tcpServer: rw.Read failed %v\n", e)
						break
					}
				}
				return
			}()
		}
	}()
	return l.Addr().String()
}

func TestCloseIdle(t *testing.T) {

	SetLogger(&logger{})

	addr := startTcpServer()
	tcpaddr, _ := net.ResolveTCPAddr("tcp", addr)

	limit := 10
	timeout := 3 * time.Second
	{
		server := &mongoServer{
			Addr:         addr,
			ResolvedAddr: tcpaddr.String(),
			tcpaddr:      tcpaddr,
			sync:         make(chan bool, 100),
			info:         &defaultServerInfo,
		}

		sockets := make([]*mongoSocket, 10)
		for i := range sockets {
			socket, _, _ := server.AcquireSocket(limit, timeout)
			assert.Nil(t, socket.dead)
			sockets[i] = socket
		}
		server.Close()
		for _, socket := range sockets {
			assert.NotNil(t, socket.dead)
			_, err := socket.conn.Write([]byte("should failed"))
			assert.Contains(t, "use of closed network connection", err.Error())
		}
	}

	{
		server := &mongoServer{
			Addr:         addr,
			ResolvedAddr: tcpaddr.String(),
			tcpaddr:      tcpaddr,
			sync:         make(chan bool, 100),
			info:         &defaultServerInfo,
		}

		sockets := make([]*mongoSocket, 10)
		for i := range sockets {
			socket, _, _ := server.AcquireSocket(limit, timeout)
			assert.Nil(t, socket.dead)
			sockets[i] = socket
		}

		for i, socket := range sockets {
			if i%2 == 0 {
				socket.Release()
			}
		}
		server.CloseIdle()

		for i, socket := range sockets {
			if i%2 != 0 {
				assert.Nil(t, socket.dead)
				_, err := socket.conn.Write([]byte("should not failed"))
				assert.NoError(t, err)
				socket.Release()
			}
			assert.NotNil(t, socket.dead)
			_, err := socket.conn.Write([]byte("should failed"))
			assert.Contains(t, "use of closed network connection", err.Error())
		}
	}
}
