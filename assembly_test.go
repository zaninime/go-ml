package ml_test

import (
	"bytes"
	"math"
	"math/rand"
	"time"

	. "github.com/zaninime/go-ml"
)

func shuffle(a []int) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

var _ = Describe("Assembly", func() {
	Describe("creation", func() {
		It("should copy the buffer in the right position", func() {
			pkt := Packet{ContentTotalSize: 4, ContentOffset: 2, Content: []byte{0xaa, 0xbb}}
			a := NewPacketAssembly(pkt)
			Expect(bytes.Equal(a.Buffer[2:], []byte{0xaa, 0xbb})).To(BeTrue())
			Expect(a.Ready()).To(BeFalse())
		})

		It("should work correctly with unary fragments", func() {
			pkt := Packet{ContentTotalSize: 4, ContentOffset: 0, Content: []byte{0xaa, 0xbb, 0xcc, 0xdd}}
			a := NewPacketAssembly(pkt)
			Expect(a.Ready()).To(BeTrue())
		})
	})

	Describe("fragment push", func() {
		pkt1 := Packet{ContentTotalSize: 4, ContentOffset: 2, Content: []byte{0xaa, 0xbb}}
		a := NewPacketAssembly(pkt1)
		It("shouldn't be ready", func() {
			Expect(a.Ready()).To(BeFalse())
		})
		It("should copy the buffer in the right position", func() {
			pkt2 := Packet{ContentTotalSize: 4, ContentOffset: 0, Content: []byte{0xcc, 0xdd}}
			a.Push(pkt2)
			Expect(bytes.Equal(a.Buffer[2:], []byte{0xaa, 0xbb})).To(BeTrue())
			Expect(bytes.Equal(a.Buffer[:2], []byte{0xcc, 0xdd})).To(BeTrue())
		})
		It("should be ready", func() {
			Expect(a.Ready()).To(BeTrue())
		})
	})

	Context("fuzzing", func() {
		rand.Seed(time.Now().UnixNano())
		nPackets := 2 + rand.Intn(10)
		dataLength := 20 + rand.Intn(65516)
		data := make([]byte, dataLength)
		dataPerPacket := int(math.Ceil(float64(dataLength) / float64(nPackets)))
		rand.Read(data)
		pkts := make([]Packet, nPackets)
		data2 := data[:]
		for i := range pkts {
			pkts[i].ContentTotalSize = uint32(dataLength)
			pkts[i].ContentOffset = uint32(i * dataPerPacket)
			if len(data2) < dataPerPacket {
				pkts[i].Content = data2
			} else {
				pkts[i].Content = data2[:dataPerPacket]
				data2 = data2[dataPerPacket:]
			}
		}
		order := make([]int, nPackets)
		for i := range order {
			order[i] = i
		}
		shuffle(order)
		a := NewPacketAssembly(pkts[order[0]])
		order = order[1:]
		It("shouldn't be ready right after creation", func() {
			Expect(a.Ready()).To(BeFalse())
		})

		It("shouldn't be ready for intermediate packets", func() {
			for _, v := range order[:len(order)-1] {
				a.Push(pkts[v])
				Expect(a.Ready()).To(BeFalse())
			}
		})

		It("should be ready after injecting the final packet", func() {
			a.Push(pkts[order[len(order)-1]])
			Expect(a.Ready()).To(BeTrue())
		})

		It("should have correctly re-assembled the content", func() {
			Expect(bytes.Equal(a.Buffer, data)).To(BeTrue())
		})
	})
})
