package test

import (
	"testing"

	"github.com/cebecifaruk/go-sscp"
	"github.com/cebecifaruk/go-sscp/test/mockconn"
)

func TestGetChannelStatistics(t *testing.T) {
	req := mockconn.ExpectSend([]byte{
		0x01, 0x03, 0x10, 0x00, 0x04, 0xD7, 0x12, 0x90, 0x6A,
	})

	res := mockconn.ExpectRecv([]byte{
		0x01, 0x83, 0x10, 0x00, 0x23, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	testSession := mockconn.NewTestConn(t, req, res)

	conn := sscp.NewPLCConnectionFrom(testSession, 1, false)

	stat, err := conn.GetChannelStatistics("channel")

	if err != nil {
		t.Error(err)
	}

	if stat.StatisticsVersion != 1 {
		t.Fatalf("Invalid Statistics Version")
	}

	if stat.SentPackets != 0 {
		t.Fatalf("Invalid Sent Packets")
	}

	if stat.RecvPackets != 0 {
		t.Fatalf("Invalid Recv Packets")
	}

	if stat.WrongPackets != 0 {
		t.Fatalf("Invalid Wrong Packets")
	}

	if stat.SentBytes != 0 {
		t.Fatalf("Invalid Sent Bytes")
	}

	if stat.RecvBytes != 0 {
		t.Fatalf("Invalid Recv Bytes")
	}

	if len(stat.Endpoints) != 1 {
		t.Fatalf("Invalid endpoints length")
	}

	if stat.Endpoints[0].AvarageCycleTime != 0 {
		t.Fatalf("Invalid Avarage Cycle Time")
	}

	if stat.Endpoints[0].MaximalCycleTime != 0 {
		t.Fatalf("Invalid Maximal Cycle Time")

	}

	if stat.Endpoints[0].MinimalCycleTime != 0 {
		t.Fatalf("Invalid Minimal Cycle Time")
	}
}
