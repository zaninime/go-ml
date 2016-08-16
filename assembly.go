package ml

import "container/list"

type assemblyHole struct {
	First, Last int
}

// PacketAssembly buffers incoming fragments and implement packet reassembly
// using the same techniques described in RFC 815
type PacketAssembly struct {
	holes  *list.List
	Buffer []byte // Contains the reassembled data when Ready()
}

// NewPacketAssembly returns a PacketAssembly initialized using the parameters
// found in the header of pkt. Next packets feeded using Push will not modify
// the total size boundary.
func NewPacketAssembly(pkt *Packet) *PacketAssembly {
	pa := PacketAssembly{}
	pa.Buffer = make([]byte, pkt.ContentTotalSize)
	pa.holes = list.New()
	fragmentFirst := int(pkt.ContentOffset)
	fragmentLast := int(pkt.ContentOffset) + len(pkt.Content) - 1
	if fragmentFirst > 0 {
		pa.holes.PushBack(assemblyHole{0, fragmentFirst - 1})
	}
	if fragmentLast < int(pkt.ContentTotalSize)-1 {
		pa.holes.PushBack(assemblyHole{fragmentLast + 1, int(pkt.ContentTotalSize) - 1})
	}
	copy(pa.Buffer[fragmentFirst:], pkt.Content)
	return &pa
}

// Push adds the packet to the assembly by copying the content into the
// buffer at the right position.
// Out of bounds packet are ignored.
func (pa *PacketAssembly) Push(pkt *Packet) {
	fragmentFirst := int(pkt.ContentOffset)
	fragmentLast := int(pkt.ContentOffset) + len(pkt.Content) - 1
	for ihole := pa.holes.Front(); ihole != nil; {
		hole := ihole.Value.(assemblyHole)
		if fragmentFirst > hole.Last {
			ihole = ihole.Next()
			continue
		}
		if fragmentLast < hole.First {
			ihole = ihole.Next()
			continue
		}
		// "delete hole"
		if fragmentFirst > hole.First {
			pa.holes.InsertBefore(assemblyHole{hole.First, fragmentFirst - 1}, ihole)
		}
		if fragmentLast < hole.Last {
			pa.holes.InsertAfter(assemblyHole{fragmentLast + 1, hole.Last}, ihole)
		}
		toRemove := ihole
		ihole = ihole.Next()
		pa.holes.Remove(toRemove)
		copy(pa.Buffer[fragmentFirst:], pkt.Content)
	}
}

// Ready returns true when the message is complete and Buffer contains
// a full copy of the message.
func (pa *PacketAssembly) Ready() bool {
	return pa.holes.Len() == 0
}
