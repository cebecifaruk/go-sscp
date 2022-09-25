package test

import (
	"testing"

	"github.com/cebecifaruk/go-sscp"
	"github.com/cebecifaruk/go-sscp/test/mockconn"
)

func TestTimeSetupExtended(t *testing.T) {
	req := mockconn.ExpectSend([]byte{
		0x01, 0x06, 0x04, 0x00, 0x02, 0x01, 0x00,
	})

	res := mockconn.ExpectRecv([]byte{
		0x01, 0x86, 0x04, 0x00, 0x08, 0x08, 0xD4, 0x40, 0x7E, 0x93, 0x41, 0xC9,
		0xAA,
	})

	testSession := mockconn.NewTestConn(t, req, res)

	conn := sscp.NewPLCConnectionFrom(testSession, 1, false)

	_, err := conn.TimeSetupExtended(sscp.RTC_GET_UTC, nil)

	if err != nil {
		t.Error(err)
	}

	// if res != 1 {
	// 	t.Fatalf("Invalid Statistics Version")
	// }

}
