package publisher

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/aprln/vwap-engine/model"
)

type Sender interface {
	Send(msg []byte) error
}

func SetUp() Publisher {
	return New(newStdoutSender())
}

func New(s Sender) Publisher {
	return Publisher{
		sender: s,
	}
}

type Publisher struct {
	sender Sender
}

func (p Publisher) GoPublish(ch <-chan model.VWAP, wg *sync.WaitGroup) {
	go p.publishForever(ch, wg)
}

func (p Publisher) publishForever(ch <-chan model.VWAP, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		vwap, more := <-ch
		if !more {
			log.Println("no more to read from the vwap channel")

			break
		}

		jsonMsg, err := json.Marshal(vwap)
		if err != nil {
			log.Printf("JSON marshal error in publisher: %v", err)

			break
		}

		if err := p.sender.Send(jsonMsg); err != nil {
			log.Printf("sender error %v", err)

			break
		}
	}
}
