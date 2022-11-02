package publisher

import "fmt"

func newStdoutSender() stdoutSender {
	return stdoutSender{}
}

type stdoutSender struct {
}

func (p stdoutSender) Send(msg []byte) error {
	fmt.Println(string(msg))

	return nil
}
