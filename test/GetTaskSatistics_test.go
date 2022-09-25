package test

import (
	"testing"

	"github.com/cebecifaruk/go-sscp"
	"github.com/cebecifaruk/go-sscp/test/mockconn"
)

func TestGetTaskStatistics(t *testing.T) {
	req := mockconn.ExpectSend([]byte{
		0x01, 0x03, 0x01, 0x00, 0x01, 0x00,
	})
	res := mockconn.ExpectRecv([]byte{
		0x01, 0x83, 0x01, 0x00, 0x32, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
		0x47, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xAD, 0xB0, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0xC1, 0xE5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0xAD, 0xB0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xA9, 0x80, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	testSession := mockconn.NewTestConn(t, req, res)

	conn := sscp.NewPLCConnectionFrom(testSession, 1, false)

	stat, err := conn.GetTaskStatistics(0)

	if err != nil {
		t.Error(err)
	}

	if stat.StatisticsVersion != 2 {
		t.Fatalf("Invalid Statistics Version")
	}

	if stat.CycleCount != 280327 {
		t.Fatalf("Invalid Cycle Count")
	}

	if stat.LastCycleDuration != 110000 {
		t.Fatalf("Invalid Last Cycle Count")
	}

	if stat.MinimalCycleDuration != 110000 {
		t.Fatalf("Invalid Minimal Cycle Duration")
	}

	if stat.MaximalCycleDuration != 240000 {
		t.Fatalf("Invalid Maximal Cycle Duration")
	}

	if stat.WaitingForDebugger != false {
		t.Fatalf("Invalid Waiting For Debugger state")
	}

	if stat.DebuggerActualUID != 0 {
		t.Fatalf("Invalid Debugger Actual UID")
	}

	if stat.DebuggerActualOffset != 0 {
		t.Fatalf("Invalid Debugger Actual Offset")
	}
}
