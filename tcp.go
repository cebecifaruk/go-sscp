package sscp

import (
	"encoding/binary"
	"net"
)

type TCPChannel struct {
	conn net.Conn
}

func NewTCPChannel(addr string) (*TCPChannel, error) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return &TCPChannel{
		conn: conn,
	}, nil
}

func (self *TCPChannel) sendFrame(frame Frame) error {
	packet := make([]byte, 5+len(frame.Payload))

	packet[0] = byte(frame.Addr)

	binary.BigEndian.PutUint16(packet[1:], 0x3FFF&frame.FunctionId)
	binary.BigEndian.PutUint16(packet[3:], uint16(len(frame.Payload)))

	copy(packet[5:], frame.Payload)

	// Send all packet over the connection
	remaining := len(packet)

	for remaining > 0 {
		n, err := self.conn.Write(packet[len(packet)-remaining:])
		if err != nil {
			return err
		}
		remaining = remaining - n
	}

	return nil
}

func (self *TCPChannel) recvFrame() (*Frame, error) {
	frame := Frame{}

	// Read header
	buf, err := self.recv(5)

	if err != nil {
		return nil, err
	}

	frame.Addr = buf[0]
	frame.FunctionId = binary.BigEndian.Uint16(buf[1:])

	// Read Payload
	frame.Payload, err = self.recv(uint(binary.BigEndian.Uint16(buf[3:])))

	if err != nil {
		return nil, err
	}

	return &frame, nil
}

func (self *TCPChannel) recv(n uint) ([]byte, error) {
	result := make([]byte, n)
	remaining := n
	for remaining > 0 {
		buf := make([]byte, remaining)

		read, err := self.conn.Read(buf)

		if err != nil {
			return nil, err
		}

		copy(result[n-remaining:], buf[0:read])
		remaining = remaining - uint(read)

	}
	return result, nil
}
