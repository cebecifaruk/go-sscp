package sscp

import (
	"encoding/binary"
	"hash/fnv"
)

type PLCStatistics struct {
	Version uint8
}

func (self *PlcConnection) GetPLCStatistics() (*PLCStatistics, error) {
	err := self.sendFrame(0x0300, []byte{})

	if err != nil {
		return nil, err
	}

	_, err = self.recvFrame(0x0300)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Get Task Statistics

func (self *PlcConnection) GetTaskStatistics(taskId uint8) error {
	self.sendFrame(0x0301, []byte{byte(taskId)})

	// TODO: Parse response

	return nil
}

// Get Channel Statistics
type Endpoint struct {
	AvarageCycleTime uint32
	MaximalCycleTime uint32
	MinimalCycleTime uint32
}

type ChannelStatistics struct {
	Version      uint8
	SentPackets  uint32
	RecvPackets  uint32
	WrongPackets uint32
	SentBytes    uint32
	RecvBytes    uint32
	Endpoints    []Endpoint
}

func (self *PlcConnection) GetChannelStatistics(channelName string) (*ChannelStatistics, error) {
	hash := fnv.New32()
	hash.Write([]byte(channelName))

	err := self.sendFrame(0x0310, hash.Sum(nil))

	if err != nil {
		return nil, err
	}

	res, err := self.recvFrame(0x0310)

	if err != nil {
		return nil, err
	}

	var endpoints []Endpoint

	for i := uint32(0); i < binary.BigEndian.Uint32(res[21:21+2]); i++ {
		offset := 22 + 12*i
		endpoints = append(endpoints, Endpoint{
			AvarageCycleTime: binary.BigEndian.Uint32(res[offset+0 : offset+4]),
			MaximalCycleTime: binary.BigEndian.Uint32(res[offset+4 : offset+8]),
			MinimalCycleTime: binary.BigEndian.Uint32(res[offset+8 : offset+12]),
		})
	}

	return &ChannelStatistics{
		Version:      res[0],
		SentPackets:  binary.BigEndian.Uint32(res[1 : 1+4]),
		RecvPackets:  binary.BigEndian.Uint32(res[5 : 5+4]),
		WrongPackets: binary.BigEndian.Uint32(res[9 : 9+4]),
		SentBytes:    binary.BigEndian.Uint32(res[13 : 13+4]),
		RecvBytes:    binary.BigEndian.Uint32(res[17 : 17+4]),
		Endpoints:    endpoints,
	}, nil
}
