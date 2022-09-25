package sscp

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

// The device identifications are described in section 6.2 of the specification
var platforms map[uint32]string = map[uint32]string{
	0x00010000: "Generic Windows",
	0x00020000: "Generic Linux",
	0x00020001: "iPLC 510",
	0x00020002: "markMX",
	0x00020003: "iPLC P-100",
	0x00020004: "mark220",
	0x00020005: "mark320",
	0x00020006: "iPLC 520",
	0x00020007: "RaspberryPi",
	0x00020008: "esg001",
	0x00020009: "mark325",
	0x0002000A: "UniPi",
	0x0002000B: "UniPi (RPi2)",
	0x00030000: "uPLC100",
	0x00030001: "M007",
	0x00030002: "HT-1",
	0x00030003: "mark150s",
	0x00030004: "imio100",
	0x00030005: "mark120",
	0x00030006: "mark125",
	0x00030007: "mark150/485s",
	0x00030008: "mark150",
	0x00030009: "mark150/485",
	0x0003000A: "mark100",
	0x0003000B: "imio105",
	0x0003000C: "icio200",
	0x0003000D: "icio205",
}

// This type represents PLC Information
type PLCInfo struct {
	SizeOfConfigBlock uint16
	SerialNumber      []byte
	Endianness        byte
	PlatformId        uint32
	RuntimeVersion    []byte
	Name              *string
	SlaveAddr         *uint8
	TCPPort           *uint16
	SSLTCPPort        *uint16
}

func (self PLCInfo) GetPlatformString() string {
	str, ok := platforms[self.PlatformId]
	if !ok {
		return ""
	}
	return str
}

// This function simply gets basic information of the PLC.
// (This functionality defined on the section 5.4.1 of the specification)
func (self *PLCConnection) GetBasicInfo(_serialnumber string, _username string, _password string) (*PLCInfo, error) {
	serialnumber := []byte(_serialnumber)
	username := []byte(_username)
	password := []byte(_password)

	serialnumber_len := byte(len(serialnumber))
	username_len := byte(len(username))
	password_len := byte(len(password))

	if password_len > 0 {
		h := md5.New()
		h.Write([]byte(password))
		password = h.Sum(nil)
		password_len = byte(len(password))
	}

	if serialnumber_len > 255 {
		return nil, fmt.Errorf("Too long serialnumber")
	}

	if username_len > 255 {
		return nil, fmt.Errorf("Too long username")
	}

	req := make([]byte, 8+serialnumber_len+username_len+password_len)

	var offset uint16 = 0

	// Version
	req[offset] = 1
	offset += 1

	// Serial Number

	req[offset] = serialnumber_len
	offset += 1
	if serialnumber_len > 0 {
		copy(req[2:], serialnumber)
		offset += uint16(serialnumber_len)
	}

	// Username

	req[offset] = username_len
	offset += 1
	if username_len > 0 {
		copy(req[offset:], username)
		offset += uint16(username_len)
	}

	// Password

	req[offset] = password_len
	offset += 1
	if password_len > 0 {
		copy(req[offset:], password)
		offset += uint16(password_len)
	}

	binary.BigEndian.PutUint16(req[offset:], 0x0000)
	offset += 2
	binary.BigEndian.PutUint16(req[offset:], 0x0000)
	offset += 2

	res, err := self.makeRequest(0x0000, req)

	if err != nil {
		return nil, err
	}

	info := PLCInfo{}

	offset = 0
	info.SizeOfConfigBlock = binary.BigEndian.Uint16(res[offset:])
	offset += 2
	serialNumberLength := res[offset]
	offset += 1
	info.SerialNumber = res[offset : byte(offset)+serialNumberLength]
	offset += uint16(serialNumberLength)
	info.Endianness = res[offset]
	offset += 1
	info.PlatformId = binary.BigEndian.Uint32(res[offset:])
	offset += 4
	runtimeVersionLength := res[offset]
	offset += 1
	info.RuntimeVersion = res[offset : offset+uint16(runtimeVersionLength)]
	offset += uint16(runtimeVersionLength)

	// Parse Tags

	if res[offset] != 0x3E {
		return nil, fmt.Errorf("Invalid Open Tag")
	}
	offset += 1

Loop:
	for {
		if offset >= uint16(len(res)) {
			return nil, fmt.Errorf("Expected Close Tag")
		}

		tagId := res[offset]
		offset += 1

		switch tagId {
		case 0x01:
			nameBuffer := []uint16{}
			for {
				c := binary.BigEndian.Uint16(res[offset:])
				offset += 2
				if c == 0x0000 {
					break
				}
				nameBuffer = append(nameBuffer, c)
			}
			name := string(utf16.Decode(nameBuffer))
			info.Name = &name
		case 0x02:
			info.SlaveAddr = &res[offset]
			offset += 1
		case 0x04:
			tcpPort := binary.BigEndian.Uint16(res[offset:])
			info.TCPPort = &tcpPort
			offset += 2
		case 0x05:
			sslTcpPort := binary.BigEndian.Uint16(res[offset:])
			info.SSLTCPPort = &sslTcpPort
			offset += 2
		case 0x3F:
			break Loop
		default:
			return nil, fmt.Errorf("Invalid Tag Id")
		}
	}

	return &info, nil
}
