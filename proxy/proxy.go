package tcpproxy

import (
	"fmt"
	"net"
)

var uniqId int = 1

type TCPProxy struct {
	Dest string
	/** buffered transport data, and flush in callback */
	Through ThroughFunc
	ln      *net.Listener
}

type ThroughFunc func(b []byte, flush ThroughFlush, tc *TCPConn) error
type ThroughFlush func(newbyte []byte) error

func (p *TCPProxy) Listen(port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Listen port %s failed with %s", port, err)
		return err
	}

	p.ln = &ln
	fmt.Printf("Listen in port %s\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connect error", err)
		} else {
			go p.handleConnection(&conn)
		}
	}
}

// hanlde incoming connection
func (p *TCPProxy) handleConnection(conn *net.Conn) {
	destConn, err := net.Dial("tcp", p.Dest)
	if err != nil {
		fmt.Printf("Connect to dest %s failed with %s\n", p.Dest, err)
		return
	}

	closeSignal := make(chan bool)
	id := fmt.Sprintf("%03d", uniqId)
	uniqId++

	fmt.Printf("[%s] New connection from %s\n", id, (*conn).RemoteAddr().String())

	defer (*conn).Close()
	defer destConn.Close()

	go p.pipe(&TCPConn{
		Type:   Send,
		Id:     id,
		Reader: conn,
		Writer: &destConn,
	}, closeSignal)

	go p.pipe(&TCPConn{
		Type:   Receive,
		Id:     id,
		Reader: &destConn,
		Writer: conn,
	}, closeSignal)

	// wait for close
	<-closeSignal

	fmt.Printf("[%s] Closed connection\n", id)
}

func (p *TCPProxy) pipe(tc *TCPConn, closeSignal chan bool) {
	for {
		// alloc 16kb
		buff, err := tc.read(0xffff)
		if err != nil {
			closeSignal <- true
			return
		}

		if p.Through != nil {
			err = p.Through(buff, func(b []byte) error {
				return tc.flush(b)
			}, tc)
		} else {
			err = tc.flush(buff)
		}

		if err != nil {
			closeSignal <- true
			return
		}
	}
}
