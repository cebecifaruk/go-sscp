package sscp

import (
	"encoding/binary"
	"fmt"
)

// This file contains SSCP driver for golang.
// A frame is consists of these fields and all frames use big endian format:
//
//    +---------+-------------+---------------+------------+---------+---------+
//    | Header  | Address(1)  | FunctionId(2) | Length(2)  | Data(n) | Footer  |
//    +---------+-------------+---------------+------------+---------+---------+
//
// For TCP Connections Header and Footer is absent.
//
//
// FunctionID bit6 = isError bit7 = isResponse (for first byte !)
// 00 -> req 01 -> Okey Res 10 -> Invalid 11 -> Error Res

type Frame struct {
	Addr       uint8
	FunctionId uint16
	Payload    []byte
}

type Channel interface {
	recvFrame() (*Frame, error)
	sendFrame(Frame) error
}

type PLCConnection struct {
	ch   Channel
	addr uint8
}

type Variable struct {
	Uid    uint32
	Offset uint32
	Length uint32
	Value  []byte
}

func NewPLCConnecetion(host string, addr uint8, reconnect bool) (*PLCConnection, error) {
	ch, err := NewTCPChannel(host)

	if err != nil {
		return nil, err
	}

	return &PLCConnection{
		ch:   ch,
		addr: addr,
	}, nil
}

func NewPLCConnecetionFromChannel(ch Channel, addr uint8, reconnect bool) PLCConnection {
	return PLCConnection{
		ch:   ch,
		addr: addr,
	}
}

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

func (self *PLCConnection) makeRequest(functionId uint16, reqPayload []byte) ([]byte, error) {
	reqFrame := Frame{
		Addr:       self.addr,
		FunctionId: functionId,
		Payload:    reqPayload,
	}

	err := self.ch.sendFrame(reqFrame)

	if err != nil {
		return nil, err
	}

	resFrame, err := self.ch.recvFrame()

	// Check address

	if resFrame.Addr != byte(self.addr) {
		return nil, fmt.Errorf("Invalid device addrss recieved")
	}

	// Check Errors and Function Code

	if resFrame.FunctionId == 0xFFFF {
		return nil, fmt.Errorf("Insufficient rights")
	}

	if resFrame.FunctionId == 0xFFFE {
		return nil, fmt.Errorf("Invalid function")
	}

	if resFrame.FunctionId == 0xFFFD {
		return nil, fmt.Errorf("Invalid protocol version")
	}

	if resFrame.FunctionId&uint16(0x3FFF) != functionId {
		return nil, fmt.Errorf("Invalid response function code")
	}

	// Check is an error response

	if (reqFrame.FunctionId & 0xC000) == 0xC0 {
		errCode := binary.BigEndian.Uint32(reqFrame.Payload)

		errString, ok := errorCodeTable[errCode]

		if ok {
			return resFrame.Payload, fmt.Errorf(errString)
			// TODO: Optional data for 0x0108, 0x0103, 0x010A, 0x0112, 0x0115, 0x0113
		}

		return resFrame.Payload, fmt.Errorf("Error response")
	}

	return resFrame.Payload, nil
}
