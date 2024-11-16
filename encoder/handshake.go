package encoder

import "bytes"

const PROTOCOL_IDENTIFIER = "BitTorrent protocol"

func MakeHandshakeMessage(infoHash string, peerID string, extended bool) []byte {
	// The handshake message is a 68-byte message that is used to establish a connection between two peers.
	// The message is structured as follows:
	// 1. The first byte is the length of the protocol identifier.
	// 2. The next 19 bytes are the protocol identifier string "BitTorrent protocol".
	// 3. The next 8 bytes are reserved for future use.
	// 4. The next 20 bytes are the SHA1 hash of the info dictionary from the .torrent file.
	// 5. The last 20 bytes are the peer ID of the sender.
	// The total length of the handshake message is 68 bytes.
	// The handshake message is used to establish a connection between two peers.
	buff := new(bytes.Buffer)
	// 1. The first byte is the length of the protocol identifier.
	buff.WriteByte(byte(len(PROTOCOL_IDENTIFIER)))
	// 2. The next 19 bytes are the protocol identifier string "BitTorrent protocol".
	buff.WriteString(PROTOCOL_IDENTIFIER)
	if extended {
		// 3. To signal support for extensions, a client must set the 20th bit from the right (counting starts at 0) in the reserved bytes to 1.
		reserved := make([]byte, 8)
		// Set the 20th bit from the right to 1
		reserved[5] = 1 << 4
		for _, b := range reserved {
			buff.WriteByte(b)
		}
	} else {
		// 3. The next 8 bytes are reserved for future use.
		for i := 0; i < 8; i++ {
			buff.WriteByte(0)
		}
	}
	// 4. The next 20 bytes are the SHA1 hash of the info dictionary from the .torrent file.
	buff.WriteString(infoHash)
	// 5. The last 20 bytes are the peer ID of the sender.
	buff.WriteString(peerID)
	return buff.Bytes()
}
