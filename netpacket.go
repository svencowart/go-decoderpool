package go_decoderpool

import (
	"bytes"
	"net"
)

// NetPacket are generated when a udp listener receives a packet
type NetPacket struct {
	Address net.Addr
	Buffer  *bytes.Buffer
	NBytes  int
}

// NewNetPacket creates a new NetPacket
func NewNetPacket(buffer []byte) NetPacket {
	buf := bytes.NewBuffer(buffer)

	return NetPacket{
		Buffer: buf,
	}
}
