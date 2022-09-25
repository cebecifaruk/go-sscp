package test

import (
	"bytes"
	"testing"

	"github.com/cebecifaruk/go-sscp"
	"github.com/cebecifaruk/go-sscp/test/mockconn"
)

func TestGetBasicInfo(t *testing.T) {
	req := mockconn.ExpectSend([]byte{
		0x00, 0x00, // Function Get Basic Info
		0x00, 0x1D, // Data length (29 Bytes)
		0x01,                         // Version
		0x00,                         // Serial number is not known
		0x05,                         // User name length
		0x61, 0x64, 0x6D, 0x69, 0x6E, // User name in utf-8 (“admin”)
		0x10, // Password length (16 Bytes)
		0x03, 0x8c, 0x0d, 0xc8, 0xa9, 0x58, 0xff, 0xea,
		0x17, 0xaf, 0x04, 0x72, 0x44, 0xfb, 0x69, 0x60, // MD5 hash for password
		0x00, 0x00, // Offset in memoryregion
		0x00, 0x00, // Size of transfered block (0 = detect only)
	})
	res := mockconn.ExpectRecv([]byte{
		0x80, 0x00, // Respond to Get Basic Info function
		0x00, 0x28, // Data length (40 bytes)
		0x04, 0x3D, // Size of whole config block (1085)
		0x08,                                           // Serial number length (8 Bytes)
		0x00, 0x00, 0x00, 0x0A, 0x14, 0xBE, 0x14, 0xB0, // Serial number
		0x00,                   // Target endianness (little)
		0x00, 0x03, 0x00, 0x07, // Platform ID (uPLC - mark150/485s)
		0x04,                   // Runtime version length (4 Bytes)
		0x22, 0xF2, 0xC0, 0x02, // Runtime version (1.0.2309.49154)
		0x3E,                                           // Open tag
		0x01,                                           // Device name tag
		0x00, 0x50, 0x00, 0x4C, 0x00, 0x43, 0x00, 0x00, // “PLC” in UTF-16
		0x02, // SSCP slave address tag
		0x01,
		0x04,       // SSCP TCP slave port tag
		0x30, 0x3A, // 12346
		0x05, // SSCP SSL slave port tag
		0x00, 0x00,
		0x3F,
	})

	testSession := mockconn.NewTestConn(t, req, res)

	conn := sscp.NewPLCConnectionFrom(testSession, 1, false)

	info, err := conn.GetBasicInfo("", "admin", "rw")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if info.SizeOfConfigBlock != 1085 {
		t.Fatalf("Invalid size of config block expected 1085 found %d", info.SizeOfConfigBlock)
	}

	expectedSerialNumber := []byte{0x00, 0x00, 0x00, 0x0A, 0x14, 0xBE, 0x14, 0xB0}
	if bytes.Compare(info.SerialNumber, expectedSerialNumber) != 0 {
		t.Fatalf("Invalid serial number expected \n%+v found \n%+v", expectedSerialNumber, info.SerialNumber)
	}

	if info.Endianness != 0 {
		t.Fatalf("Invalid endianness expected \n%+v found \n%+v", 0, info.Endianness)
	}

	if info.PlatformId != 0x00030007 {
		t.Fatalf("Invalid plaform id expected \n% X found \n% X", 0x00030007, info.PlatformId)
	}

	if bytes.Compare(info.RuntimeVersion, []byte{0x22, 0xF2, 0xC0, 0x02}) != 0 {
		t.Fatalf("Invalid runtime version expected \n% X found \n% X", []byte{0x22, 0xF2, 0xC0, 0x02}, info.RuntimeVersion)
	}

	if !(info.SSLTCPPort != nil && *info.SSLTCPPort == 0) {
		t.Fatalf("Invalid SSLTCPPort")
	}

	if !(info.SlaveAddr != nil && *info.SlaveAddr == 1) {
		t.Fatalf("Invalid SlaveAddr")
	}

	if !(info.Name != nil && *info.Name == "PLC") {
		if info.Name != nil {
			t.Fatalf("Invalid name expected \n%+v found \n%+v", "PLC", info.Name)
		}
	}
}
