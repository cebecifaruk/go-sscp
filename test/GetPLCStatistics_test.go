package test

import (
	"testing"

	"github.com/cebecifaruk/go-sscp"
	"github.com/cebecifaruk/go-sscp/test/mockconn"
)

func TestGetPLCStatistics(t *testing.T) {
	req := mockconn.ExpectSend([]byte{
		0x01, 0x03, 0x00, 0x00, 0x00,
	})

	res := mockconn.ExpectRecv([]byte{
		0x01,       //SSCP Address
		0x83, 0x00, // Get PLC Statistics Response
		0x00, 0x73, //Data Length (115)
		0x04,       // Statistics Version
		0x00,       // Runtime Stat Block Id
		0x01,       // Block Version
		0x00, 0x1C, // Block Length
		0x01,                                           // Normal Task Count
		0x00,                                           // Max Task Id
		0x01,                                           // Evaluator State
		0x00,                                           // Run Mode (Full Run)
		0x00, 0x00, 0x00, 0x00, 0x4f, 0x6d, 0x40, 0x80, // Uptime
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // Running tasks mask
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Tasks with exception mask
		0x01,       // Memory Stat block Id
		0x01,       // Block Version
		0x00, 0x10, // Block Size (16)
		0x20, 0x8f, // Total heap size (8335 kB)
		0x1e, 0x5f, // Free heap after runtime startup (7775 kB)
		0x1d, 0x73, // Free heap after image load and process (7539 kB)
		0x01, 0xff, // Total space available for image (511 kB)
		0x01, 0x23, // Free space after image save (291 kB)
		0x00, 0x01, // Retain size (1 kB)
		0x02, 0x00, // Allocator total size (512 kB)
		0x02, 0x00, // Allocator free space (512 kB)
		0x02,       // Sections statistics block ID
		0x01,       // Block version
		0x00, 0x06, // Block size
		0x00, 0x8E, // VMEX section used (142 kB)
		0x00, 0x40, // RTCM section used (64 kB)
		0x00, 0x0D, // Other sections used (13 kB)
		0x03,       // RCware DB statistics block ID
		0x01,       // Block version
		0x00, 0x15, // Block size (21 Bytes)
		0x00,                   // Client status (Disabled)
		0x00, 0x00, 0x00, 0x00, // Records saved
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Last save time
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Last request time
		0x04,       // Proxy statistics block ID
		0x01,       // Block version
		0x00, 0x17, // Block size
		0x00, // Proxy status (Disabled)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Proxy ID
		0x00, // Slots total
		0x00, // Slots free
	})

	testSession := mockconn.NewTestConn(t, req, res)

	conn := sscp.NewPLCConnectionFrom(testSession, 1, false)

	stat, err := conn.GetPLCStatistics()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if stat.StatisticsVersion != 4 {
		t.Fatalf("Invalid StatisticsVersion")
	}

	if stat.NormalTaskCount != 1 {
		t.Fatalf("Invalid NormalTaskCount")
	}

	if stat.MaxTaskId != 0 {
		t.Fatalf("Invalid MaxTaskId")
	}

	if stat.EvaluatorState != 1 {
		t.Fatalf("Invalid EvaluatorState")
	}

	if stat.RunMode != 0 {
		t.Fatalf("Invalid RunMode")
	}

	if stat.UpTime != 0x000000004f6d4080 {
		t.Fatalf("Invalid UpTime")
	}

	if stat.RunningTasks != 0x0000000000000001 {
		t.Fatalf("Invalid RunningTasks")
	}

	if stat.TasksWithException != 0x0000000000000000 {
		t.Fatalf("Invalid TasksWithException")
	}

	if stat.TotalHeap != 8335 {
		t.Fatalf("Invalid TotalHeap")
	}

	if stat.FreeHeapBeforeLoad != 7775 {
		t.Fatalf("Invalid FreeHeap")
	}

	if stat.FreeHeap != 7539 {
		t.Fatalf("Invalid FreeHeapBeforeLoad")
	}

	if stat.TotalCodeSpace != 511 {
		t.Fatalf("Invalid TotalCodeSpace")
	}

	if stat.FreeCodeSpace != 291 {
		t.Fatalf("Invalid FreeCodeSpace")
	}

	if stat.RetainSize != 1 {
		t.Fatalf("Invalid RetainSize")
	}

	if stat.AllocatorTotalSize != 512 {
		t.Fatalf("Invalid AllocatorTotalSize")
	}

	if stat.AllocatorFreeSpace != 512 {
		t.Fatalf("Invalid AllocatorFreeSpace")
	}

	if stat.VMEXSection != 142 {
		t.Fatalf("Invalid VMEXSection")
	}

	if stat.RTCMSection != 64 {
		t.Fatalf("Invalid RTCMSection")
	}

	if stat.OtherSections != 13 {
		t.Fatalf("Invalid OtherSections")
	}
}
