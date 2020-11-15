package go_decoderpool

type Pool struct {
	concurrency int
	done        chan error
	NetPackets  chan NetPacket
	Messages    chan Message
}

func (p *Pool) Run() {
	for i := 0; i < p.concurrency; i++ {
		decoder := NewDecoder(p.done, p.NetPackets, p.Messages)

		go decoder.Run()
	}

	for {
		select {
		case <-p.done:
			close(p.NetPackets)
			close(p.Messages)

			return
		}
	}
}

func NewPool(done chan error, concurrency int) Pool {
	return Pool{
		concurrency: concurrency,
		done:        done,
		NetPackets:  make(chan NetPacket, 100000),
		Messages:    make(chan Message),
	}
}
