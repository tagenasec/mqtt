package packets

import (
	"bytes"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
)

func TestPubrecEncode(t *testing.T) {
	require.Contains(t, expectedPackets, Pubrec)
	for i, wanted := range expectedPackets[Pubrec] {

		if !encodeTestOK(wanted) {
			continue
		}

		require.Equal(t, uint8(5), Pubrec, "Incorrect Packet Type [i:%d] %s", i, wanted.desc)
		pk := new(PubrecPacket)
		copier.Copy(pk, wanted.packet.(*PubrecPacket))

		require.Equal(t, Pubrec, pk.Type, "Mismatched Packet Type [i:%d] %s", i, wanted.desc)
		require.Equal(t, Pubrec, pk.FixedHeader.Type, "Mismatched FixedHeader Type [i:%d] %s", i, wanted.desc)

		buf := new(bytes.Buffer)
		err := pk.Encode(buf)
		require.NoError(t, err, "Expected no error writing buffer [i:%d] %s", i, wanted.desc)
		encoded := buf.Bytes()

		require.Equal(t, len(wanted.rawBytes), len(encoded), "Mismatched packet length [i:%d] %s", i, wanted.desc)
		require.Equal(t, byte(Pubrec<<4), encoded[0], "Mismatched fixed header packets [i:%d] %s", i, wanted.desc)
		require.EqualValues(t, wanted.rawBytes, encoded, "Mismatched byte values [i:%d] %s", i, wanted.desc)

		require.Equal(t, wanted.packet.(*PubrecPacket).PacketID, pk.PacketID, "Mismatched Packet ID [i:%d] %s", i, wanted.desc)

	}
}

func BenchmarkPubrecEncode(b *testing.B) {
	pk := new(PubrecPacket)
	copier.Copy(pk, expectedPackets[Pubrec][0].packet.(*PubrecPacket))

	buf := new(bytes.Buffer)
	for n := 0; n < b.N; n++ {
		pk.Encode(buf)
	}
}

func TestPubrecDecode(t *testing.T) {
	require.Contains(t, expectedPackets, Pubrec)
	for i, wanted := range expectedPackets[Pubrec] {

		if !decodeTestOK(wanted) {
			continue
		}

		require.Equal(t, uint8(5), Pubrec, "Incorrect Packet Type [i:%d] %s", i, wanted.desc)

		pk := newPacket(Pubrec).(*PubrecPacket)
		err := pk.Decode(wanted.rawBytes[2:]) // Unpack skips fixedheader.

		if wanted.failFirst != nil {
			require.Error(t, err, "Expected error unpacking buffer [i:%d] %s", i, wanted.desc)
			require.Equal(t, wanted.failFirst, err.Error(), "Expected fail state; %v [i:%d] %s", err.Error(), i, wanted.desc)
			continue
		}

		require.NoError(t, err, "Error unpacking buffer [i:%d] %s", i, wanted.desc)

		require.Equal(t, wanted.packet.(*PubrecPacket).PacketID, pk.PacketID, "Mismatched Packet ID [i:%d] %s", i, wanted.desc)
	}

}

func BenchmarkPubrecDecode(b *testing.B) {
	pk := newPacket(Pubrec).(*PubrecPacket)
	pk.FixedHeader.decode(expectedPackets[Pubrec][0].rawBytes[0])

	for n := 0; n < b.N; n++ {
		pk.Decode(expectedPackets[Pubrec][0].rawBytes[2:])
	}
}

func TestPubrecValidate(t *testing.T) {
	pk := newPacket(Pubrec).(*PubrecPacket)
	pk.FixedHeader.decode(expectedPackets[Pubrec][0].rawBytes[0])

	b, err := pk.Validate()
	require.NoError(t, err)
	require.Equal(t, Accepted, b)

}

func BenchmarkPubrecValidate(b *testing.B) {
	pk := newPacket(Pubrec).(*PubrecPacket)
	pk.FixedHeader.decode(expectedPackets[Pubrec][0].rawBytes[0])

	for n := 0; n < b.N; n++ {
		pk.Validate()
	}
}
