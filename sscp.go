package sscp

import (
	"encoding/binary"
	"fmt"
	"net"
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

type PLCConnection struct {
	conn net.Conn
	addr uint8
}

func NewPLCConnection(host string, addr uint8, reconnect bool) (*PLCConnection, error) {
	conn, err := net.Dial("tcp", host)

	if err != nil {
		return nil, err
	}

	return &PLCConnection{
		conn: conn,
		addr: addr,
	}, nil
}

func NewPLCConnectionFrom(conn net.Conn, addr uint8, reconnect bool) PLCConnection {
	return PLCConnection{
		conn: conn,
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

func (self *PLCConnection) sendFrame(frame Frame) error {
	packet := make([]byte, 5+len(frame.Payload))

	packet[0] = byte(frame.Addr)

	binary.BigEndian.PutUint16(packet[1:], 0x3FFF&frame.FunctionId)
	binary.BigEndian.PutUint16(packet[3:], uint16(len(frame.Payload)))

	copy(packet[5:], frame.Payload)

	err := self.send(packet)

	if err != nil {
		return err
	}

	return nil
}

func (self *PLCConnection) sendNonAddrFrame(frame Frame) error {
	packet := make([]byte, 4+len(frame.Payload))

	binary.BigEndian.PutUint16(packet[0:], 0x3FFF&frame.FunctionId)
	binary.BigEndian.PutUint16(packet[2:], uint16(len(frame.Payload)))

	copy(packet[4:], frame.Payload)

	err := self.send(packet)

	if err != nil {
		return err
	}

	return nil
}

func (self *PLCConnection) recvFrame() (*Frame, error) {
	frame := Frame{}

	// Read header
	buf, err := self.recv(5)

	if err != nil {
		return nil, err
	}

	frame.Addr = buf[0]
	frame.FunctionId = binary.BigEndian.Uint16(buf[1:])

	// Read Payload
	frame.Payload, err = self.recv(uint(binary.BigEndian.Uint16(buf[3:])))

	if err != nil {
		return nil, err
	}

	return &frame, nil
}

func (self *PLCConnection) recvNonAddrFrame() (*Frame, error) {
	frame := Frame{}

	// Read header
	buf, err := self.recv(4)

	if err != nil {
		return nil, err
	}

	frame.FunctionId = binary.BigEndian.Uint16(buf[0:])

	// Read Payload
	frame.Payload, err = self.recv(uint(binary.BigEndian.Uint16(buf[2:])))

	if err != nil {
		return nil, err
	}

	return &frame, nil
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

// Send all packet over the connection
func (self *PLCConnection) send(packet []byte) error {
	remaining := len(packet)

	for remaining > 0 {
		n, err := self.conn.Write(packet[len(packet)-remaining:])
		if err != nil {
			return err
		}
		remaining = remaining - n
	}

	return nil
}

func (self *PLCConnection) makeRequest(functionId uint16, reqPayload []byte) ([]byte, error) {
	reqFrame := Frame{
		Addr:       self.addr,
		FunctionId: functionId,
		Payload:    reqPayload,
	}

	if functionId == 0x0000 {
		err := self.sendNonAddrFrame(reqFrame)
		if err != nil {
			return nil, err
		}
	} else {
		err := self.sendFrame(reqFrame)
		if err != nil {
			return nil, err
		}
	}

	var resFrame *Frame
	var err error

	if functionId == 0x0101 {
		return nil, nil
	} else if functionId == 0x0000 {
		resFrame, err = self.recvNonAddrFrame()
		if err != nil {
			return nil, err
		}
	} else {
		resFrame, err = self.recvFrame()
		if err != nil {
			return nil, err
		}

		// Check address

		if resFrame.Addr != byte(self.addr) {
			return nil, fmt.Errorf("Invalid device addrss recieved")
		}
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

	if (resFrame.FunctionId & 0xC000) == 0xC0 {
		errCode := binary.BigEndian.Uint32(resFrame.Payload)

		errString, ok := errorCodeTable[errCode]

		if ok {
			return nil, fmt.Errorf(errString+" with payload %+v", resFrame.Payload)
		}

		return nil, fmt.Errorf("Error response"+" with payload %+v", resFrame.Payload)
	}

	return resFrame.Payload, nil
}

func handleRuntimeError() error {
	if r := recover(); r != nil {
		return fmt.Errorf("Rutime Error %+v", r)
	}
	return nil
}
