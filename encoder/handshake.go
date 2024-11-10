package encoder

import "bytes"

func MakeHandshakeMessage(infoHash string, peerID string) []byte {
	// The handshake message is a 68-byte message that is used to establish a connection between two peers.
	// The message is structured as follows:
	// 1. The first byte is the length of the protocol identifier.
	// 2. The next 19 bytes are the protocol identifier string "BitTorrent protocol".
	// 3. The next 8 bytes are reserved for future use.
	// 4. The next 20 bytes are the SHA1 hash of the info dictionary from the .torrent file.
	// 5. The last 20 bytes are the peer ID of the sender.
	// The total length of the handshake message is 68 bytes.
	// The handshake message is used to establish a connection between two peers.
	const protocolIdentifier = "BitTorrent protocol"
	buff := new(bytes.Buffer)
	// 1. The first byte is the length of the protocol identifier.
	buff.WriteByte(byte(len(protocolIdentifier)))
	// 2. The next 19 bytes are the protocol identifier string "BitTorrent protocol".
	buff.WriteString(protocolIdentifier)
	// 3. The next 8 bytes are reserved for future use.
	for i := 0; i < 8; i++ {
		buff.WriteByte(0)
	}
	// 4. The next 20 bytes are the SHA1 hash of the info dictionary from the .torrent file.
	buff.WriteString(infoHash)
	// 5. The last 20 bytes are the peer ID of the sender.
	buff.WriteString(peerID)
	return buff.Bytes()
}
