package decoder

import (
	"encoding/binary"
	"fmt"
)

type PeerMessage struct {
	Length  uint32
	Id      uint8
	Payload []byte
}

func NewPeerMessage(id uint8, payload []byte) *PeerMessage {
	return &PeerMessage{
		Length:  uint32(1 + len(payload)),
		Id:      id,
		Payload: payload,
	}
}

func (d *Decoder) DecodePeerMessage() *PeerMessage {
	length := d.readInt32()
	if int(length) != len(d.data[d.offset:]) {
		fmt.Printf("Error: length of message is %d but the remaining data is %d", length, len(d.data[d.offset:]))
		return nil
	}
	id := d.readInt8()
	payload := d.data[d.offset:]
	return NewPeerMessage(uint8(id), payload)
}

func DecodePieceMessage(data []byte) (pieceIndex, byteOffset int, dataBlock []byte) {
	if len(data) < 13 {
		fmt.Printf("Error: piece message is too short: %d", len(data))
		return
	}
	data = data[5:] // Skip the length and ID
	pieceIndex = int(binary.BigEndian.Uint32(data[0:4]))
	byteOffset = int(binary.BigEndian.Uint32(data[4:8]))
	dataBlock = data[8:]
	return
}

func (pm *PeerMessage) Encode() []byte {
	buff := make([]byte, 5+len(pm.Payload))
	binary.BigEndian.PutUint32(buff[0:4], pm.Length)
	buff[4] = pm.Id
	copy(buff[5:], pm.Payload)
	return buff
}

func (pm *PeerMessage) String() string {
	switch pm.Id {
	case CHOKE, UNCHOKE, INTERESTED, NOT_INTERESTED:
		return MessageNames[pm.Id]
	case REQUEST:
		pieceIndex := binary.BigEndian.Uint32(pm.Payload[0:4])
		begin := binary.BigEndian.Uint32(pm.Payload[4:8])
		length := binary.BigEndian.Uint32(pm.Payload[8:12])
		return fmt.Sprintf("%s piece index: %d byte offset: %d length: %d\n", MessageNames[pm.Id], pieceIndex, begin, length)
	default:
		return MessageNames[pm.Id]
	}
}

const ( // Message types
	CHOKE          = iota // NO PAYLOAD
	UNCHOKE               // NO PAYLOAD
	INTERESTED            // NO PAYLOAD
	NOT_INTERESTED        // NO PAYLOAD
	HAVE
	BITFIELD
	REQUEST
	PIECE
	CANCEL
)

var MessageNames = map[uint8]string{
	CHOKE:          "CHOKE",
	UNCHOKE:        "UNCHOKE",
	INTERESTED:     "INTERESTED",
	NOT_INTERESTED: "NOT_INTERESTED",
	HAVE:           "HAVE",
	BITFIELD:       "BITFIELD",
	REQUEST:        "REQUEST",
	PIECE:          "PIECE",
	CANCEL:         "CANCEL",
}

func BitfieldMessage(payload []byte) *PeerMessage {
	return NewPeerMessage(BITFIELD, payload)
}

func ChokeMessage() *PeerMessage {
	return NewPeerMessage(CHOKE, []byte{})
}

func UnchokeMessage() *PeerMessage {
	return NewPeerMessage(UNCHOKE, []byte{})
}

func InterestedMessage() *PeerMessage {
	return NewPeerMessage(INTERESTED, []byte{})
}

func NotInterestedMessage() *PeerMessage {
	return NewPeerMessage(NOT_INTERESTED, []byte{})
}

func RequestMessage(index, begin, length uint32) *PeerMessage {
	buff := make([]byte, 12)
	binary.BigEndian.PutUint32(buff[0:4], index)
	binary.BigEndian.PutUint32(buff[4:8], begin)
	binary.BigEndian.PutUint32(buff[8:12], length)
	return NewPeerMessage(REQUEST, buff)
}
