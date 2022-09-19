package sscp

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

import (
	"net"
)

type PlcConnection struct {
	conn      net.Conn
	addr      uint8
	reconnect bool
}

func NewPlcConnecetion(host string, addr uint8, reconnect bool) (*PlcConnection, error) {
	conn, err := net.Dial("tcp", host)

	if err != nil {
		return nil, err
	}

	return &PlcConnection{
		conn:      conn,
		addr:      addr,
		reconnect: reconnect,
	}, nil
}
