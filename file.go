package sscp

import (
	"encoding/binary"
	"fmt"
	"time"
)

type FileInfo struct {
	Size uint32
	Time time.Time
	CRC  uint16
}

const (
	FILE_RT_IMAGE         = "/sys/sexm"
	FILE_RT_UPGRADE_IMAGE = "/sys/rt"
	FILE_SHARK_IMAGE_1    = "/sys/sex1"
	FILE_SHARK_IMAGE_2    = "/sys/sex2"
	FILE_PLC_CAPS         = "/sys/caps"
	FILE_VAR_DIRECT       = "/var/direct"
	FILE_LOGS             = "/log"
)

// This functionality defined on the section 5.6.1 of the specification
func (self *PLCConnection) InitiateDataSend(filename string, size uint32, time time.Time) error {
	encodedFilename := []byte(filename)

	if len(encodedFilename) > 64 {
		return fmt.Errorf("Filename is too long")
	}

	req := make([]byte, 13+len(encodedFilename))
	self.makeRequest(0x0200, req)
	return fmt.Errorf("Not implemented yet.")
}

// This functionality defined on the section 5.6.1 of the specification
func (self *PLCConnection) SendDataChunk(offset uint32, data []byte) error {
	return fmt.Errorf("Not implemented yet.")
}

// This functionality defined on the section 5.6.1 of the specification
func (self *PLCConnection) FinishDataSend(crc uint16) error {
	req := make([]byte, 2)
	binary.BigEndian.PutUint16(req, crc)

	// Response frame does not contain any important data.
	_, err := self.makeRequest(0x0202, req)

	return err
}

// This functionality defined on the section 5.6.2 of the specification
func (self *PLCConnection) InitiateDataReceive(filename string) (*FileInfo, error) {
	encodedFilename := []byte(filename)

	if len(encodedFilename) > 64 {
		return nil, fmt.Errorf("Filename is too long")
	}

	req := make([]byte, 1+len(encodedFilename))
	req[0] = byte(len(encodedFilename))

	copy(req[1:], encodedFilename)

	res, err := self.makeRequest(0x0210, req)

	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Size: binary.BigEndian.Uint32(res[0:4]),
		Time: FromDateTime(binary.BigEndian.Uint64(res[4:12])),
		CRC:  binary.BigEndian.Uint16(res[12:14]),
	}, err
}

// This functionality defined on the section 5.6.2 of the specification
func (self *PLCConnection) ReceiveDataChunk(offset uint32) (uint32, []byte, error) {
	req := make([]byte, 4)
	binary.BigEndian.PutUint32(req, offset)
	res, err := self.makeRequest(0x0211, req)

	if err != nil {
		return 0, nil, err
	}

	return binary.BigEndian.Uint32(res[0:4]), res[4:], nil
}

// This function sends a complete file.
func (self *PLCConnection) SendFile(filename string, content []byte) error {
	return nil
}

// This function recieves a complete file.
func (self *PLCConnection) RecvFile(filename string) ([]byte, error) {
	return nil, nil
}
