package sscp

import (
	"encoding/binary"
	"fmt"
)

type Variable struct {
	Uid    uint32
	Offset uint32
	Length uint32
	Value  []byte
}

// This functionality defined on the section 5.8.1.1 of the specification
// It simply takes a list of variables and mutates their values.
func (self *PLCConnection) ReadVariablesDirectly(vars []*Variable, taskId *uint8) error {
	if len(vars) > 64 {
		return fmt.Errorf("Too many variables")
	}

	var payload []byte

	if taskId != nil {
		payload = make([]byte, 2+12*len(vars))
	} else {
		payload = make([]byte, 1+12*len(vars))
	}

	offset := 0
	payload[offset] = 0x80

	if taskId != nil {
		payload[offset] |= 0x20
	}

	offset += 1

	if taskId != nil {
		payload[offset] = *taskId
		offset += 1
	}

	for _, v := range vars {
		binary.BigEndian.PutUint32(payload[offset:], v.Uid)
		offset += 4
		binary.BigEndian.PutUint32(payload[offset:], v.Offset)
		offset += 4
		binary.BigEndian.PutUint32(payload[offset:], v.Length)
		offset += 4
	}

	res, err := self.makeRequest(0x0500, payload)

	if err != nil {
		return err
	}

	offset = 0

	for _, v := range vars {
		if offset < len(res) {
			len := int(v.Length)
			v.Value = res[offset : offset+len]
			offset += len
		}
	}

	return nil
}

// This functionality defined on the section 5.8.1.2 of the specification
func (self *PLCConnection) WriteVariablesDirectly(vars []*Variable) error {
	numOfVars := uint32(len(vars))
	totalRawLength := uint32(0)

	if numOfVars > 255 {
		return fmt.Errorf("Too many variables")
	}

	for _, v := range vars {
		totalRawLength += v.Length
	}

	req := make([]byte, 2+12*numOfVars+totalRawLength)

	req[0] = 0x80
	req[1] = byte(numOfVars)

	for i, v := range vars {
		frame := req[2+12*i : 2+12*i+12]
		binary.BigEndian.PutUint32(frame[0:4], v.Uid)
		binary.BigEndian.PutUint32(frame[4:8], v.Offset)
		binary.BigEndian.PutUint32(frame[8:12], v.Length)
	}

	offset := 2 + 12*numOfVars

	for _, v := range vars {
		copy(req[offset:offset+v.Length], v.Value)
		offset += v.Length
	}

	_, err := self.makeRequest(0x0510, req)

	if err != nil {
		return err
	}

	return nil
}
