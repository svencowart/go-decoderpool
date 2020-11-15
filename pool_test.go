package go_decoderpool

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPool_Run(t *testing.T) {
	done := make(chan error)
	defer close(done)

	pool := NewPool(done, 1)

	go pool.Run()

	timestamp := time.Now()
	jsonTimestamp, _ := timestamp.MarshalJSON()
	jsonStr := fmt.Sprintf(`{"id":1,"timestamp":%s}`, jsonTimestamp)
	jsonBytes := []byte(jsonStr)
	packet := NewNetPacket(jsonBytes)

	pool.NetPackets <- packet

	var jsonMsg []byte
	select {
	case msg := <-pool.Messages:
		jsonMsg, _ = json.Marshal(msg)
	}

	expectedMsg := `{"id":1,"timestamp":"(.*)"}`
	t.Log(string(jsonMsg))
	assert.Regexp(t, expectedMsg, string(jsonMsg), "should match regexp")
}

func TestPool_Concurrency(t *testing.T) {
	count := 10000
	done := make(chan error)
	defer close(done)

	var elapsedCon1 time.Duration
	var elapsedCon4 time.Duration

	func() {
		pool := NewPool(done, 2)

		go pool.Run()

		start := time.Now()
		go func() {
			for i := 0; i < count; i++ {
				timestamp := time.Now()
				jsonTimestamp, _ := timestamp.MarshalJSON()
				jsonStr := fmt.Sprintf(`{"id":%v,"timestamp":%s}`, i, jsonTimestamp)
				jsonBytes := []byte(jsonStr)
				packet := NewNetPacket(jsonBytes)

				pool.NetPackets <- packet
			}
		}()

		counter := 0
		for _ = range pool.Messages {
			counter += 1
			if counter == count {
				break
			}
		}

		end := time.Now()
		elapsedCon1 = end.Sub(start)
	}()

	func() {
		pool := NewPool(done, 4)

		go pool.Run()

		start := time.Now()
		go func() {
			for i := 0; i < count; i++ {
				timestamp := time.Now()
				jsonTimestamp, _ := timestamp.MarshalJSON()
				jsonStr := fmt.Sprintf(`{"id":%v,"timestamp":%s}`, i, jsonTimestamp)
				jsonBytes := []byte(jsonStr)
				packet := NewNetPacket(jsonBytes)

				pool.NetPackets <- packet
			}
		}()

		counter := 0
		for _ = range pool.Messages {
			counter += 1
			if counter == count {
				break
			}
		}

		end := time.Now()
		elapsedCon4 = end.Sub(start)
	}()

	t.Log("pool_concurrency_1 runtime: ", elapsedCon1)
	t.Log("pool_concurrency_4 runtime: ", elapsedCon4)
	assert.LessOrEqualf(t, int64(elapsedCon4), int64(elapsedCon1), "elapsed time for a pool of 1 should be more than a pool of 4 - pool_concurrency_1: %v / pool_concurrency_4: %v", elapsedCon1, elapsedCon4)
}