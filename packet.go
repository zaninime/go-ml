package ml

import (
	"encoding/binary"
	"errors"
)

// Packet is the data structure representing a ML-library packet.
// This structure is used by parsers and serializers.
type Packet struct {
	ContentOffset    uint32 // offset of data contained in this fragment
	ContentTotalSize uint32 // content size of the transmission
	LocalID          uint32
	RemoteID         uint32
	Sequence         uint32 // multiple fragments share the same sequence number
	Type             uint8  // content type
	//FuncLength       uint16 // not implemented
	Content []byte
}

// ErrPacketTooShort is retured by ParsePacket if the byte slice is shorter
// than the header length
var ErrPacketTooShort = errors.New("ml: packet too short")

// ParsePacket returns a parsed packet from a byte slice.
func ParsePacket(data []byte) (*Packet, error) {
	const headerLength = 23
	if len(data) < headerLength {
		return nil, ErrPacketTooShort
	}
	packet := Packet{}
	packet.ContentOffset = binary.BigEndian.Uint32(data[0:4])
	packet.ContentTotalSize = binary.BigEndian.Uint32(data[4:8])
	packet.LocalID = binary.BigEndian.Uint32(data[8:12])
	packet.RemoteID = binary.BigEndian.Uint32(data[12:16])
	packet.Sequence = binary.BigEndian.Uint32(data[16:20])
	packet.Type = data[20]
	packet.Content = data[23:]

	return &packet, nil
}
