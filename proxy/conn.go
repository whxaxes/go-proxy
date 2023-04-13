package tcpproxy

import (
	"fmt"
	"net"
)

const (
	Send = iota
	Receive
)

type TCPConn struct {
	Type   int
	Id     string
	Reader *net.Conn
	Writer *net.Conn
}

func (tc *TCPConn) read(maxsize int) ([]byte, error) {
	buffer := make([]byte, maxsize)
	n, err := (*tc.Reader).Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

// flush data to writer
func (tc *TCPConn) flush(data []byte) error {
	n := len(data)

	ReaderAddr := (*tc.Reader).RemoteAddr().String()
	WriterAddr := (*tc.Writer).RemoteAddr().String()
	if tc.Type == Send {
		fmt.Printf("[%s] %s >>> %s (%d bytes)\n", tc.Id, ReaderAddr, WriterAddr, n)
	} else {
		fmt.Printf("[%s] %s <<< %s (%d bytes)\n", tc.Id, WriterAddr, ReaderAddr, n)
	}

	_, err := (*tc.Writer).Write(data)
	return err
}
