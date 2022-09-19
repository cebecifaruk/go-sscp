package sscp

import (
	"crypto/md5"
	"encoding/binary"
)

type LoginResponse struct {
	ProtoVersion uint8
	MaxDataSize  uint16
	RightGroup   uint8
	ImageGUID    uint32

	// Optional Fields
	DeviceName  *string
	SSCPAddress *uint8
	SSCPTCPPort *uint16
	SSCPSSLPort *uint16
}

func (self *PlcConnection) Login(username string, password string, proxyId string, maxDataSize uint16) (*LoginResponse, error) {
	const funcCode = 0x0100

	// Get username as byte buffer
	_username := []byte(username)
	_username_len := len(_username)

	// Get password hash as byte buffer
	h := md5.New()
	h.Write([]byte(password))
	_password := h.Sum(nil)
	_password_len := len(_password)

	_proxyId := []byte(proxyId)
	_proxyId_len := len(_proxyId)

	payload := make([]byte, 6+_username_len+_password_len+_proxyId_len)

	payload[0] = 0x07
	payload[1] = byte((maxDataSize | 0xFF00) >> 8)
	payload[2] = byte((maxDataSize | 0x00FF) >> 0)

	// Username
	payload[3] = byte(_username_len)
	if _username_len > 0 {
		copy(payload[4:], _username)
	}

	// Password
	payload[4+_username_len] = byte(_password_len)
	if _password_len > 0 {
		copy(payload[5+_username_len:], _password)
	}

	// Proxy Id
	payload[5+_username_len+_password_len] = byte(_proxyId_len)
	if _proxyId_len > 0 {
		copy(payload[6+_username_len+_password_len:], _proxyId)
	}

	self.sendFrame(funcCode, payload)

	// Recieve Response

	res, err := self.recvFrame(funcCode)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		ProtoVersion: res[0],
		MaxDataSize:  binary.BigEndian.Uint16(res[1:3]),
		RightGroup:   res[3],
		ImageGUID:    binary.BigEndian.Uint32(res[4:8]),
	}, nil
}

func (self *PlcConnection) Logout() error {
	self.sendFrame(0x0101, []byte{})

	return nil
}
