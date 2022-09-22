package sscp

import (
	"encoding/binary"
	"fmt"
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
	SerialNumber   []byte
	Endianness     byte
	Platform       string
	RuntimeVersion []byte
	Name           string
	SlaveAddr      uint8
	TCPPort        uint16
	SSLTCPPort     uint16
}

func (self *PLCInfo) parseInfoTags(buffer []byte) error {
	if buffer[0] != 0x3E {
		return fmt.Errorf("Invalid Open Tag")
	}

	if buffer[len(buffer)-1] != 0x3F {
		return fmt.Errorf("Invalid Close Tag")
	}

	tagsBuffer := buffer[1 : len(buffer)-1]

	for i := 0; i < len(tagsBuffer); i++ {
		tagId := tagsBuffer[i]
		i += 1

		switch tagId {
		case 0x01:
			for {
				i += 2
			}
		case 0x02:
			self.SlaveAddr = tagsBuffer[i]
			i += 1
		case 0x04:
			self.TCPPort = binary.BigEndian.Uint16(tagsBuffer[i : i+2])
			i += 2
		case 0x05:
			self.SSLTCPPort = binary.BigEndian.Uint16(tagsBuffer[i : i+2])
			i += 2
		default:
			return fmt.Errorf("Invalid Tag Id")
		}
	}

	return nil
}

// This function simply gets basic information of the PLC.
// (This functionality defined on the section 5.4.1 of the specification)
func (self *PLCConnection) GetBasicInfo(serialnumber string, username string, password string) (*PLCInfo, error) {
	_sn := []byte(serialnumber)
	_username := []byte(username)
	_password := []byte(password)

	len_sn := byte(len(_sn))
	len_username := byte(len(_username))
	len_password := byte(len(_password))

	if len_sn > 255 {
		return nil, fmt.Errorf("Too long serialnumber")
	}

	if len_username > 255 {
		return nil, fmt.Errorf("Too long username")
	}

	if len_password > 255 {
		return nil, fmt.Errorf("Too long password")
	}

	req := make([]byte, 8+len_username+len_password+len_sn)

	req[0] = 1

	// Serial Number

	req[1] = len_sn
	if len(_sn) > 0 {
		copy(req[2:], _sn)
	}

	// Username

	req[3+len_sn] = len_username
	if len(_sn) > 0 {
		copy(req[4+len_sn:], _username)
	}

	// Password

	req[5+len_sn+len_username] = len_password
	if len(_sn) > 0 {
		copy(req[6+len_sn+len_username:], _password)
	}

	// req[7+len_sn+len_username+len_password] = 0x00
	// req[8+len_sn+len_username+len_password] = 0x00
	// req[9+len_sn+len_username+len_password] = 0x00
	// req[10+len_sn+len_username+len_password] = 0x00

	res, err := self.makeRequest(0x0000, req)

	if err != nil {
		return nil, err
	}

	serialNumberLength := res[2]
	runtimeVersionLength := res[0]

	info := PLCInfo{
		SerialNumber:   res[3 : 3+serialNumberLength],
		Endianness:     res[4+serialNumberLength],
		Platform:       platforms[binary.BigEndian.Uint32(res[5+serialNumberLength:9+serialNumberLength])],
		RuntimeVersion: res[9+serialNumberLength : 9+serialNumberLength+runtimeVersionLength],
	}

	info.parseInfoTags(res[9+serialNumberLength:])

	return &info, nil
}
