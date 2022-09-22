package test

import (
	"testing"

	"github.com/cebecifaruk/go-sscp"
	"github.com/cebecifaruk/go-sscp/test/mockconn"
)

func TestWriteVariablesInDirectMode(t *testing.T) {

	req := mockconn.ExpectSend([]byte{
		0x01, 0x05, 0x10, 0x00, 0x1D, 0x80, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x02, 0x35,
	})

	res := mockconn.ExpectRecv([]byte{
		0x01, 0x85, 0x10, 0x00, 0x00,
	})

	testSession := mockconn.NewTestConn(t, req, res)

	conn := sscp.NewPLCConnectionFrom(testSession, 1, false)

	v1 := sscp.Variable{
		Uid:    0x00000001,
		Offset: 0x00000000,
		Length: 0x00000001,
		Value:  []byte{0x01},
	}

	v2 := sscp.Variable{
		Uid:    0x00000002,
		Offset: 0x00000000,
		Length: 0x00000002,
		Value:  []byte{0x02, 0x35},
	}

	err := conn.WriteVariablesDirectly([]*sscp.Variable{&v1, &v2})

	if err != nil {
		t.Error(err)
	}
}

// TODO: Create Write in File Mode Test
