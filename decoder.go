package go_decoderpool

import (
	"encoding/json"
	"time"
)

type Decoder struct {
	done   chan error
	input  <-chan NetPacket
	output chan<- Message
}

func (d *Decoder) Run() {
	for packet := range d.input {
		msg, err := decode(packet)
		if err != nil {
			d.done <- err
		}

		d.output <- msg
	}
}

func decode(packet NetPacket) (Message, error) {
	var msg Message

	time.Sleep(20 * time.Millisecond)

	err := json.Unmarshal(packet.Buffer.Bytes(), &msg)
	if err != nil {
		//log.Println(err)
		return msg, err
	}

	return msg, nil
}

func NewDecoder(done chan error, input <-chan NetPacket, output chan<- Message) Decoder {
	return Decoder{
		done:   done,
		input:  input,
		output: output,
	}
}
