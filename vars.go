package sscp

import (
	"encoding/binary"
	"fmt"
)

// This functionality defined on the section 5.8.1.1 of the specification
func (self *PLCConnection) ReadVariablesDirectly(vars []*Variable) error {
	if len(vars) > 64 {
		return fmt.Errorf("Too many variables")
	}

	payload := make([]byte, 1+12*len(vars))
	payload[0] = 0x80

	for i, v := range vars {
		binary.BigEndian.PutUint32(payload[1+3*i:], v.Uid)
		binary.BigEndian.PutUint32(payload[5+3*i:], v.Offset)
		binary.BigEndian.PutUint32(payload[9+3*i:], v.Length)
	}

	err := self.sendFrame(0x0500, payload)

	if err != nil {
		return err
	}

	res, err := self.recvFrame(0x0500)

	if err != nil {
		return err
	}

	var ptr uint32 = 0

	for _, v := range vars {
		v.Value = res[ptr : ptr+v.Length]
		ptr += v.Length
	}

	return nil
}

// This functionality defined on the section 5.8.1.2 of the specification
func (self *PLCConnection) WriteVariableDirectly() {}
