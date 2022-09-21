package mockconn

import "net"

type MockAddr struct{}

func (self MockAddr) Network() string { return "" }
func (self MockAddr) String() string  { return "" }

func NewMockAddr() net.Addr {
	return MockAddr{}
}
