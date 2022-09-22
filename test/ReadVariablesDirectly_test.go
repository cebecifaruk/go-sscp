package test

import (
	"bytes"
	"testing"

	"github.com/cebecifaruk/sscp"
	"github.com/cebecifaruk/sscp/test/mockconn"
)

func TestReadVariablesInDirectMode(t *testing.T) {

	req := mockconn.ExpectSend([]byte{
		0x01, 0x05, 0x00, 0x00, 0x25, 0x80, 0x00, 0x00, 0x22, 0xBE, 0x00, 0x00, 0x00,
		0xD9, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x22, 0xC0, 0x00, 0x00, 0x00, 0xDA,
		0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x22, 0xBF, 0x00, 0x00, 0x01, 0x84, 0x00,
		0x00, 0x00, 0x04,
	})

	res := mockconn.ExpectRecv([]byte{
		0x01, 0x85, 0x00, 0x00, 0x07, 0x00, 0x00, 0x02, 0x42, 0x48, 0x00, 0x00,
	})

	testSession := mockconn.NewTestConn(t, req, res)

	conn := sscp.NewPLCConnectionFrom(testSession, 1, false)

	v1 := sscp.Variable{
		Uid:    0x000022BE,
		Offset: 0x000000D9,
		Length: 0x00000001,
	}

	v2 := sscp.Variable{
		Uid:    0x000022C0,
		Offset: 0x000000DA,
		Length: 0x00000002,
	}

	v3 := sscp.Variable{
		Uid:    0x000022BF,
		Offset: 0x00000184,
		Length: 0x00000004,
	}

	err := conn.ReadVariablesDirectly([]*sscp.Variable{&v1, &v2, &v3})

	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(v1.Value, []byte{0x00}) != 0 {
		t.Fatalf("Invalid Value for variable 1")
	}

	if bytes.Compare(v2.Value, []byte{0x00, 0x02}) != 0 {
		t.Fatalf("Invalid Value for variable 1")
	}

	if bytes.Compare(v3.Value, []byte{0x42, 0x48, 0x00, 0x00}) != 0 {
		t.Fatalf("Invalid Value for variable 1")
	}
}

// TODO: Create Read in File Mode Test
