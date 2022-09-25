package sscp

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
)

const (
	PROXY_DISABLED        = 0
	PROXY_NOTUSED         = 1
	PROXY_IDLE            = 2
	PROXY_CONNECTED       = 3
	PROXY_UNAUTHORIZED    = 4
	PROXY_NOTAVAILABLE    = 5
	PROXY_FAILEDTOCONNECT = 6
	PROXY_HOSTNOTFOUND    = 7
	PROXY_CONNECTING      = 8
	PROXY_PAGENOTFOUND    = 9
	PROXY_DBERROR         = 10
)

const (
	RUN_MODE_FULL_RUN                    = 0
	RUN_MODE_CommunicationOnly           = 1
	RUN_MODE_EvaluationOnly              = 2
	RUN_MODE_Commissioning               = 3
	RUN_MODE_CommunicationsWithTransform = 4
	RUN_MODE_PrepareOnly                 = 5
	RUN_MODE_StartDisabledBySwitch       = 32
	RUN_MODE_InvalidImageVersion         = 33
	RUN_MODE_NoMemoryForImage            = 34
)

const (
	EVAL_STATE_Stopped                     = 0
	EVAL_STATE_RunningNormalTasks          = 1
	EVAL_STATE_StoppingExecution           = 2
	EVAL_STATE_RunningExceptionStateTask   = 3
	EVAL_STATE_ExceptionStateTaskFailed    = 4
	EVAL_STATE_NoExceptionStateTaskDefined = 5
	EVAL_STATE_Commissioning               = 6
	EVAL_STATE_InvalidImage                = 7
	EVAL_STATE_NoImage                     = 8
	EVAL_STATE_WaitingForDebugger          = 9
	EVAL_STATE_PreparedForStart            = 10
)

type PLCStatistics struct {
	StatisticsVersion  uint8
	BlockLength        uint16
	NormalTaskCount    uint8
	MaxTaskId          uint8
	EvaluatorState     uint8
	RunMode            uint8
	UpTime             uint64
	RunningTasks       uint64
	TasksWithException uint64
	TotalHeap          uint16
	FreeHeapBeforeLoad uint16
	FreeHeap           uint16
	TotalCodeSpace     uint16
	FreeCodeSpace      uint16
	RetainSize         uint16
	AllocatorTotalSize uint16
	AllocatorFreeSpace uint16
	VMEXSection        uint16
	RTCMSection        uint16
	OtherSections      uint16
	ClientStatus       uint8
	RecordsSaved       uint32
	LastSaveTime       uint64
	LastRequestTime    uint64
	ProxyStatus        uint8
	ProxyId            []byte
	SlotsTotal         uint8
	SlotsFree          uint8
}

type Endpoint struct {
	AvarageCycleTime uint32
	MaximalCycleTime uint32
	MinimalCycleTime uint32
}

type ChannelStatistics struct {
	StatisticsVersion uint8
	SentPackets       uint32
	RecvPackets       uint32
	WrongPackets      uint32
	SentBytes         uint32
	RecvBytes         uint32
	Endpoints         []Endpoint
}

type TaskStatistics struct {
	StatisticsVersion    uint8
	CycleCount           uint64
	LastCycleDuration    uint64
	AvarageCycleDuration uint64
	MinimalCycleDuration uint64
	MaximalCycleDuration uint64
	WaitingForDebugger   bool
	DebuggerActualUID    uint32
	DebuggerActualOffset uint32
}

// This functionality defined on the section 5.7.1 of the specification
func (self *PLCConnection) GetPLCStatistics() (*PLCStatistics, error) {
	_, err := self.makeRequest(0x0300, []byte{})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// This functionality defined on the section 5.7.2 of the specification
func (self *PLCConnection) GetTaskStatistics(taskId uint8) (*TaskStatistics, error) {
	res, err := self.makeRequest(0x0301, []byte{byte(taskId)})

	if err != nil {
		return nil, err
	}

	if len(res) < 41 {
		return nil, fmt.Errorf("Invalid response, length of response is not enough to parse")
	}

	result := TaskStatistics{
		StatisticsVersion:    res[0],
		CycleCount:           binary.BigEndian.Uint64(res[1:]),
		LastCycleDuration:    binary.BigEndian.Uint64(res[9:]),
		AvarageCycleDuration: binary.BigEndian.Uint64(res[17:]),
		MinimalCycleDuration: binary.BigEndian.Uint64(res[25:]),
		MaximalCycleDuration: binary.BigEndian.Uint64(res[33:]),
	}

	if result.StatisticsVersion == 2 {
		result.WaitingForDebugger = res[41] != 0
		result.DebuggerActualUID = binary.BigEndian.Uint32(res[42:])
		result.DebuggerActualOffset = binary.BigEndian.Uint32(res[46:])
	}

	return &result, nil
}

// This functionality defined on the section 5.7.3 of the specification
func (self *PLCConnection) GetChannelStatistics(channelName string) (*ChannelStatistics, error) {
	hash := fnv.New32()
	hash.Write([]byte(channelName))

	res, err := self.makeRequest(0x0310, hash.Sum(nil))

	if err != nil {
		return nil, err
	}

	var endpoints []Endpoint

	for i := uint16(0); i < binary.BigEndian.Uint16(res[21:]); i++ {
		offset := 23 + 12*i
		endpoints = append(endpoints, Endpoint{
			AvarageCycleTime: binary.BigEndian.Uint32(res[offset+0 : offset+4]),
			MaximalCycleTime: binary.BigEndian.Uint32(res[offset+4 : offset+8]),
			MinimalCycleTime: binary.BigEndian.Uint32(res[offset+8 : offset+12]),
		})
	}

	return &ChannelStatistics{
		StatisticsVersion: res[0],
		SentPackets:       binary.BigEndian.Uint32(res[1:]),
		RecvPackets:       binary.BigEndian.Uint32(res[5:]),
		WrongPackets:      binary.BigEndian.Uint32(res[9:]),
		SentBytes:         binary.BigEndian.Uint32(res[13:]),
		RecvBytes:         binary.BigEndian.Uint32(res[17:]),
		Endpoints:         endpoints,
	}, nil
}
