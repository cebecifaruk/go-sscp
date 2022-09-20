package sscp

import (
	"encoding/binary"
	"fmt"
)

// Error Code Table
var errorCodeTable map[uint32]string = map[uint32]string{
	0x0000: "No Error",
	0x0001: "No Response",
	0x0002: "Failed To Connect",
	0x0003: "Not Implemented",
	0x0004: "Invalid Function Received",
	0x0101: "Wrong Login",
	0x0102: "No Such File",
	0x0103: "No Such Variable",
	0x0104: "No Such Task",
	0x0105: "Wrong Order",
	0x0106: "Wrong Parameter",
	0x0107: "Invalid Group Id",
	0x0108: "Transmission In Progress",
	0x0109: "Not Registered",
	0x010A: "Write Failed",
	0x010B: "Not All Data Received",
	0x010C: "Invalid Crc",
	0x010D: "Data Too Long",
	0x010E: "Too Long Use File Transfer",
	0x010F: "File Name Too Long",
	0x0110: "Variable Count Limit Exceed",
	0x0111: "Out Of Bounds",
	0x0112: "Size Mismatch",
	0x0113: "Operation Denied",
	0x0114: "Not Logged",
	0x0115: "Invalid State",
	0x0116: "Unknown Channel",
	0x0117: "Driver Command Timeout",
	0x0118: "Unknown Driver Command",
	0x0119: "No Resources Available",
	0x011A: "Chunk Read Failed",
	0x011B: "Chunk Write Failed",
	0x011C: "No Such Metadata",
	0x011D: "Async",
	0x0801: "SysCmd_NewImage",
	0x0802: "SysCmd_InvalidImageArea",
	0x0803: "SysCmd_CreateBootImage",
	0x0804: "SysCmd_WarmReboot",
	0x0805: "SysCmd_ColdReboot",
	0x0806: "SysCmd_StartPlc",
	0x0807: "SysCmd_StopPlc",
	0x0808: "SysCmd_SetMacAddress",
	0x0809: "SysCmd_Timeout",
	0x080C: "SysCmdRequestActive",
	0x080D: "SysCmdWaitTimeout",
	0x080A: "Already Running",
	0x080B: "Already Stopped",
}

func (self *PLCConnection) recv(n uint) ([]byte, error) {
	result := make([]byte, n)
	remaining := n
	for remaining > 0 {
		buf := make([]byte, remaining)

		read, err := self.conn.Read(buf)

		if err != nil {
			return nil, err
		}

		copy(result[n-remaining:], buf[0:read])
		remaining = remaining - uint(read)

	}
	return result, nil
}

func (self *PLCConnection) send(data []byte) error {
	remaining := len(data)
	for remaining > 0 {
		n, err := self.conn.Write(data[len(data)-remaining:])
		if err != nil {
			return err
		}
		remaining = remaining - n
	}
	return nil
}

func (self *PLCConnection) sendFrame(functionId uint16, payload []byte) error {
	functionId = 0x3FFF & functionId
	payloadLen := len(payload)

	packet := []byte{
		byte(self.addr),
		byte((functionId & 0xFF00) >> 8),
		byte(functionId & 0x00FF),
		byte((payloadLen & 0xFF00) >> 8),
		byte(payloadLen & 0x00FF),
	}

	packet = append(packet, payload...)

	err := self.send(packet)

	if err != nil {
		return err
	}

	return nil
}

func (self *PLCConnection) recvFrame(functionId uint16) ([]byte, error) {
	functionId = 0x3FFF & functionId

	// Read header
	buf, err := self.recv(5)

	if err != nil {
		return nil, err
	}

	// Read Payload
	payload, err := self.recv(uint(buf[3])<<8 | uint(buf[4]))

	if err != nil {
		return nil, err
	}

	// Check address

	if buf[0] != byte(self.addr) {
		return nil, fmt.Errorf("Invalid device addrss recieved")
	}

	// Check Errors and Function Code

	functionCode := uint16(buf[1])<<8 | uint16(buf[2])

	if functionCode == 0xFFFF {
		return nil, fmt.Errorf("Insufficient rights")
	}

	if functionCode == 0xFFFE {
		return nil, fmt.Errorf("Invalid function")
	}

	if functionCode == 0xFFFD {
		return nil, fmt.Errorf("Invalid protocol version")
	}

	if functionCode&uint16(0x3FFF) != functionId {
		return nil, fmt.Errorf("Invalid response function code")
	}

	// Check is an error response

	if (buf[1] & 0x40) > 0 {
		errCode := binary.BigEndian.Uint32(payload[0:4])

		errString, ok := errorCodeTable[errCode]

		if ok {
			return payload, fmt.Errorf(errString)
		}

		return payload, fmt.Errorf("Error response")
	}

	return payload, nil
}
