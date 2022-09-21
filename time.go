package sscp

import (
	"encoding/binary"
	"time"
)

// TODO: Not implemented
// Converts time data type to 100ns timestamp from 01-01-0001
// https://msdn.microsoft.com/en-us/library/system.datetime.ticks(v=vs.110).aspx
func toDateTime(t time.Time) uint64 {
	return 0
}

// TODO: Not implemented
// Converts 100ns timestamp from 01-01-0001 to time data type
// https://msdn.microsoft.com/en-us/library/system.datetime.ticks(v=vs.110).aspx
func fromDateTime(t uint64) time.Time {
	// Convert to us
	t = t / 10
	// Change time origin to 01-01-1970
	t = t - 621_355_968_000_000_000
	return time.Now()
}

// This functionality defined on the section 5.9.1 of the specification
func (self *PLCConnection) TimeSetup(t *time.Time) (*time.Time, error) {
	req := []byte{}

	if t != nil {
		req = make([]byte, 4)
		binary.BigEndian.PutUint64(req, toDateTime(*t))
	}

	res, err := self.makeRequest(0x0602, req)

	if err != nil {
		return nil, err
	}

	result := fromDateTime(binary.BigEndian.Uint64(res))

	return &result, nil
}

const (
	RTC_GET_UTC             = 0x01
	RTC_GET_LOCAL           = 0x02
	RTC_SET_UTC             = 0x10
	RTC_SET_LOCAL           = 0x11
	RTC_GET_TIMEZONE_OFFSET = 0x20
	RTC_GET_DAYLIGHT_OFFSET = 0x21
)

// TODO: Not implemented
func (self *PLCConnection) TimeSetupExtended(t *time.Time) (*time.Time, error) {
	return nil, nil
}
