package publisher

func NewMockSender() MockSender {
	return MockSender{
		msgChan: make(chan []byte, 1),
	}
}

type MockSender struct {
	msgChan chan []byte
}

func (m MockSender) Send(msg []byte) error {
	m.msgChan <- msg

	return nil
}

func (m MockSender) Read() []byte {
	return <-m.msgChan
}

func (m MockSender) Close() {
	close(m.msgChan)
}
