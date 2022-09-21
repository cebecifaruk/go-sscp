package mockconn

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

type TestConn struct {
	t   *testing.T
	ops []ExpectedOperation
	i   int
	ptr int
}

func NewTestConn(t *testing.T, expectedOps ...ExpectedOperation) net.Conn {
	conn := TestConn{
		t:   t,
		ops: expectedOps,
		i:   0,
		ptr: 0,
	}

	return &conn
}

func (self *TestConn) Read(b []byte) (n int, err error) {
	if self.i >= len(self.ops) {
		self.t.Fatalf("Expected no operation but found RECV")
	}
	op := self.ops[self.i]

	if !op.IsRecv() {
		self.t.Fatalf(
			"Expected channel operation is %s with payload \n% x \nbut found RECV",
			op.Name(),
			op.Payload,
		)
	}

	l := copy(b, op.Payload[self.ptr:])
	self.ptr += l
	if self.ptr > len(op.Payload) {
		self.i += 1
		self.ptr = 0
	}
	return l, nil
}

func (self *TestConn) Write(b []byte) (n int, err error) {
	if self.i >= len(self.ops) {
		self.t.Fatalf("Expected no operation but found SEND")
	}
	op := self.ops[self.i]
	self.i += 1

	fmt.Println(op.IsSend())

	if !op.IsSend() {
		self.t.Fatalf(
			"Expected channel operation is %s with payload \n% x \nbut found SEND with payload: \n% x",
			op.Name(),
			op.Payload,
			b,
		)
	}

	if bytes.Compare(op.Payload, b) != 0 {
		self.t.Fatalf(
			"Expected payload is \n% x \nbut found \n% x",
			op.Payload,
			b,
		)
	}

	return len(b), nil
}

func (self *TestConn) Close() error {
	if self.i >= len(self.ops) {
		self.t.Fatalf("Expected no operation but found CLOSE")
	}
	op := self.ops[self.i]
	self.i += 1

	if !op.IsClose() {
		self.t.Fatalf(
			"Expected channel operation is %s with payload \n% x \nbut found CLOSE",
			op.Name(),
			op.Payload,
		)
	}

	return nil
}

func (self *TestConn) LocalAddr() net.Addr {
	return NewMockAddr()
}

func (self *TestConn) RemoteAddr() net.Addr {
	return NewMockAddr()
}

func (self *TestConn) SetDeadline(t time.Time) error {
	return nil
}

func (self *TestConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (self *TestConn) SetWriteDeadline(t time.Time) error {
	return nil
}
