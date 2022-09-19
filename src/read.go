package sscp

import "fmt"

type Varaible struct {
	Uid    uint32
	Offset uint32
	Length uint32
	Value  []byte
}

func (self *PlcConnection) ReadVariables(flags uint8, vars []Varaible) error {
	if len(vars) > 64 {
		return fmt.Errorf("Too many variables")
	}

	payload := make([]byte, 1+12*len(vars))
	payload[0] = 0x80

	for i, v := range vars {
		// Uid
		payload[1+3*i] = byte((v.Uid & 0xFF000000) >> 24)
		payload[2+3*i] = byte((v.Uid & 0x00FF0000) >> 16)
		payload[3+3*i] = byte((v.Uid & 0x0000FF00) >> 8)
		payload[4+3*i] = byte((v.Uid & 0x000000FF) >> 0)

		// Offset
		payload[5+3*i] = byte((v.Offset & 0xFF000000) >> 24)
		payload[6+3*i] = byte((v.Offset & 0x00FF0000) >> 16)
		payload[7+3*i] = byte((v.Offset & 0x0000FF00) >> 8)
		payload[8+3*i] = byte((v.Offset & 0x000000FF) >> 0)

		// Length
		payload[9+3*i] = byte((v.Length & 0xFF000000) >> 24)
		payload[10+3*i] = byte((v.Length & 0x00FF0000) >> 16)
		payload[11+3*i] = byte((v.Length & 0x0000FF00) >> 8)
		payload[12+3*i] = byte((v.Length & 0x000000FF) >> 0)
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
