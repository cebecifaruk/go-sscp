package test

import (
	"bytes"
	"sscp"
	"sscp/test/mockconn"
	"testing"
)

func TestLogin(t *testing.T) {
	req := mockconn.ExpectSend([]byte{
		0x01, 0x01, 0x00, 0x00, 0x1B, 0x07, 0x28, 0x00, 0x05, 0x61, 0x64, 0x6D, 0x69,
		0x6E, 0x10, 0x03, 0x8c, 0x0d, 0xc8, 0xa9, 0x58, 0xff, 0xea, 0x17, 0xaf, 0x04,
		0x72, 0x44, 0xfb, 0x69, 0x60, 0x00,
	})
	res := mockconn.ExpectRecv([]byte{
		0x01, 0x81, 0x00, 0x00, 0x1B, 0x07, 0x00, 0xE4, 0xFF, 0xF0, 0x2A, 0x9D, 0x0B,
		0x2A, 0x37, 0x75, 0x44, 0xB6, 0xAF, 0x28, 0x21, 0x05, 0xA2, 0xCA, 0x00, 0x3E,
		0x03, 0x58, 0x45, 0x44, 0xF8, 0x3F,
	})

	testSession := mockconn.NewTestConn(t, req, res, mockconn.ExpectClose())

	conn := sscp.NewPLCConnecetionFromConnection(testSession, 1, false)

	result, err := conn.Login("admin", "rw", "", 10240)

	if err != nil {
		t.Error(err)
	}

	if result.ProtoVersion != 7 {
		t.Fatalf("Invalid Proto Version: Expected 7 found %d", result.ProtoVersion)
	}

	if result.MaxDataSize != 228 {
		t.Fatalf("Invalid Max Data Size: Expected 228 found %d", result.MaxDataSize)
	}

	if result.RightGroup != 0xFF {
		t.Fatalf("Invalid Right Group: Expected 0xFF found 0x% X", result.RightGroup)
	}

	guid := []byte{0xF0, 0x2A, 0x9D, 0x0B, 0x2A, 0x37, 0x75, 0x44, 0xB6, 0xAF, 0x28, 0x21, 0x05, 0xA2, 0xCA, 0x00}
	if bytes.Compare(result.ImageGUID[:], guid) != 0 {
		t.Fatalf("Invalid Image GUID: Expected % x found % x", guid, result.ImageGUID)
	}

	// TODO: Check optionals

}