package ml_test

import (
	. "github.com/zaninime/go-ml"
)

var _ = Describe("Packet", func() {
	Describe("parsing", func() {
		It("should parse header of a valid packet successfully", func() {
			pkt := []byte{0xde, 0xad, 0xc0, 0xde, 0x00, 0x0f, 0xf1, 0xce, 0xba, 0xdd, 0xca, 0xfe, 0x1c, 0xeb, 0x00, 0xda, 0, 0, 0, 10, 0x0a, 0, 2, 0x01, 0x02}
			parsed, err := ParsePacket(pkt)

			Expect(err).NotTo(HaveOccurred())
			Expect(parsed.ContentOffset).To(BeEquivalentTo(0xdeadc0de), "Content offset")
			Expect(parsed.ContentTotalSize).To(BeEquivalentTo(0x000ff1ce), "Total content size")
			Expect(parsed.LocalID).To(BeEquivalentTo(0xbaddcafe), "Local ID")
			Expect(parsed.RemoteID).To(BeEquivalentTo(0x1ceb00da), "Remote ID")
			Expect(parsed.Sequence).To(BeEquivalentTo(10), "Sequence number")
			Expect(parsed.Type).To(BeEquivalentTo(0x0a), "Packet type")
			Expect(len(parsed.Content)).To(Equal(2), "Content length")
		})

		It("should return an error if the packet is shorter than the protocol header length", func() {
			pkt := []byte{0x01, 0x02}
			_, err := ParsePacket(pkt)

			Expect(err).To(MatchError(ErrPacketTooShort))
		})
	})
})
