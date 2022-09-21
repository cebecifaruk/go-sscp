package mockconn

const (
	CLOSE = iota
	RECV
	SEND
)

type ExpectedOperation struct {
	Operation byte
	Payload   []byte
}

func ExpectClose() ExpectedOperation {
	return ExpectedOperation{
		Operation: CLOSE,
		Payload:   nil,
	}
}

func ExpectRecv(payload []byte) ExpectedOperation {
	return ExpectedOperation{
		Operation: RECV,
		Payload:   payload,
	}
}

func ExpectSend(payload []byte) ExpectedOperation {
	return ExpectedOperation{
		Operation: SEND,
		Payload:   payload,
	}
}

func (self ExpectedOperation) IsClose() bool {
	return self.Operation == CLOSE
}

func (self ExpectedOperation) IsRecv() bool {
	return self.Operation == RECV
}

func (self ExpectedOperation) IsSend() bool {
	return self.Operation == SEND
}

func (self ExpectedOperation) Name() string {
	switch self.Operation {
	case SEND:
		return "SEND"
	case RECV:
		return "RECV"
	case CLOSE:
		return "CLOSE"
	}

	return "UNKNOWN"
}

// func NewSessionOp(name string, payload []byte) SessionOp {
// 	return SessionOp{
// 		name:    name,
// 		payload: payload,
// 	}
// }

// // type Channel interface {
// // 	recvFrame() (*Frame, error)
// // 	sendFrame(Frame) error
// // }
// // func (self *TestChannel) sendFrame(frame sscp.Frame) error {
// // 	packet := make([]byte, len(frame.Payload)+5)
// // 	packet[0] = frame.Addr

// // 	binary.BigEndian.PutUint16(packet[1:], frame.FunctionId)
// // 	binary.BigEndian.PutUint16(packet[3:], uint16(len(frame.Payload)))

// // 	if bytes.Compare(self.expectedReq, packet) == 0 {
// // 		return nil
// // 	}

// // 	return fmt.Errorf("Expected : %+v \nResult: %+v")
// // }
