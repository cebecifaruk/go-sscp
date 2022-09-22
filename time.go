package sscp

import (
	"encoding/binary"
	"time"
)

const (
	RTC_GET_UTC             = 0x01
	RTC_GET_LOCAL           = 0x02
	RTC_SET_UTC             = 0x10
	RTC_SET_LOCAL           = 0x11
	RTC_GET_TIMEZONE_OFFSET = 0x20
	RTC_GET_DAYLIGHT_OFFSET = 0x21
)

// Converts time data type to 100ns timestamp from 01-01-0001
// https://msdn.microsoft.com/en-us/library/system.datetime.ticks(v=vs.110).aspx
func ToDateTime(t time.Time) uint64 {
	return ((uint64(t.UnixMicro()) * 10) + 621355968000000000)

}

// Converts 100ns timestamp from 01-01-0001 to time data type
// https://msdn.microsoft.com/en-us/library/system.datetime.ticks(v=vs.110).aspx
func FromDateTime(t uint64) time.Time {
	return time.UnixMicro(int64((t - 621_355_968_000_000_000) / 10))
}

// This functionality defined on the section 5.9.1 of the specification
func (self *PLCConnection) TimeSetup(t *time.Time) (*time.Time, error) {
	req := []byte{}

	if t != nil {
		req = make([]byte, 8)
		binary.BigEndian.PutUint64(req, ToDateTime(*t))
	}

	res, err := self.makeRequest(0x0602, req)

	if err != nil {
		return nil, err
	}

	result := FromDateTime(binary.BigEndian.Uint64(res))

	return &result, nil
}

func (self *PLCConnection) TimeSetupExtended(command byte, t *time.Time) (*time.Time, error) {
	req := []byte{}

	if t != nil {
		req = make([]byte, 10)
		req[0] = command
		req[1] = 0
		binary.BigEndian.PutUint64(req, ToDateTime(*t))
	} else {
		req = make([]byte, 2)
		req[0] = command
		req[1] = 0
	}

	res, err := self.makeRequest(0x0604, req)

	if err != nil {
		return nil, err
	}

	result := FromDateTime(binary.BigEndian.Uint64(res))

	return &result, nil
}
