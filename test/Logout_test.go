package test

import (
	"sscp"
	"sscp/test/mockconn"
	"testing"
)

func TestLogout(t *testing.T) {
	req := mockconn.ExpectSend([]byte{
		0x01, 0x01, 0x01, 0x00, 0x00,
	})

	testSession := mockconn.NewTestConn(t, req, mockconn.ExpectClose())

	conn := sscp.NewPLCConnectionFrom(testSession, 1, false)

	err := conn.Logout()

	if err != nil {
		t.Error(err)
	}
}
