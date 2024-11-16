package command

import (
	"fmt"
	"math"
	"net"

	d "github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

const (
	BLOCK_LENGTH         = 16 * 1024
	PIECE_MESSAGE_LENGTH = BLOCK_LENGTH + 13
)

// DownloadPiece downloads a piece from a peer and and returns the piece data
func DownloadPiece(peerAddr string, torrentLength, torrentPieceLength int, torrentInfoHash, torrentPieceHash string, pieceIndex int, isLastPiece bool) ([]byte, error) {
	// Connect to the peer
	conn, err := helloPeer(torrentInfoHash, peerAddr)
	if err != nil {
		return nil, fmt.Errorf("error while handshaking with peer: %v", err)
	}

	isLastPieceLengthSameAsOthers := true
	if isLastPiece && torrentLength%torrentPieceLength != 0 {
		isLastPieceLengthSameAsOthers = false
	}

	// Num of blocks in a piece, this is the number of requests we will send for a piece
	// If the piece is the last piece and the piece does not divide evenly into the piece length, adjust the piece length
	pieceLength, numOfBlocks := 0, 0
	if isLastPiece && !isLastPieceLengthSameAsOthers {
		pieceLength = torrentLength % torrentPieceLength
		numOfBlocks = int(math.Ceil(float64(pieceLength) / float64(BLOCK_LENGTH))) // Number of blocks in the last piece
		// fmt.Printf("Adjusting piece length from: %d to: %d and the number of blocks from: %d to: %d\n", torrentPieceLength, pieceLength, pieceLength/BLOCK_LENGTH, numOfBlocks)
	} else {
		pieceLength = torrentPieceLength
		numOfBlocks = pieceLength / BLOCK_LENGTH
	}

	// Break the piece into blocks of 16 kiB (16 * 1024 bytes) and send a request message for each block
	// If the piece is the last piece and the piece does not divide evenly into the piece length, adjust the piece length
	requests := createRequests(pieceLength, numOfBlocks, pieceIndex)
	err = sendRequests(requests, conn)
	if err != nil {
		return nil, fmt.Errorf("error while sending request messages: %v", err)
	}

	// Receive the piece message from the peer
	pieceStream, err := receivePiece(conn, pieceLength)
	if err != nil {
		return nil, fmt.Errorf("error while receiving piece message: %v", err)
	}

	// Chunk the piece message into blocks of 16 kiB + 13 bytes
	blocks := chunkPieceStream(pieceStream, pieceLength, isLastPiece, isLastPieceLengthSameAsOthers)

	// Decode each piece message and reconstruct the piece
	pieceReconstructed := make([]byte, 0)
	for _, block := range blocks {
		_, _, dataBlock := d.DecodePieceMessage(block)
		// fmt.Printf("Decoded piece message for piece: %d, byte offset: %d, length: %d\n", pieceIndex, byteOffset, len(dataBlock))
		pieceReconstructed = append(pieceReconstructed, dataBlock...)
	}

	// Check if the piece reconstructed sha1 hash matches the piece hash in the torrent file
	if fmt.Sprintf("%x", utils.SHA1Hash(pieceReconstructed)) == torrentPieceHash {
		fmt.Println("Piece hash matches the piece hash in the torrent file")
		return pieceReconstructed, nil
	} else {
		return nil, fmt.Errorf("error piece hash does not match the piece hash in the torrent file, expected: %s, got: %x", torrentPieceHash, utils.SHA1Hash(pieceReconstructed))
	}
}

// Exchange multiple peer messages with a peer to ensure we can download a piece from the peer
// If the peer is ready, we return the connection to the peer
func helloPeer(torrentInfoHash string, peerAddr string) (net.Conn, error) {
	// Connect to the peer
	conn, err := Handshake(torrentInfoHash, peerAddr)
	if err != nil {
		return nil, fmt.Errorf("error while handshaking with peer: %v", err)
	}
	// fmt.Printf("Handshake successful with peer: %s\n", peerAddr)

	// Wait for the bitfield message
	buff := make([]byte, 1024)
	size, err := conn.Read(buff)
	if err != nil {
		return nil, fmt.Errorf("error while reading bitfield message: %v", err)
	}
	data := buff[:size]
	// d.LogMessage(data, true)

	// Decode the bitfield message
	decoder := d.NewDecoder(data)
	pm := decoder.DecodePeerMessage()
	if pm.Id != d.BITFIELD {
		return nil, fmt.Errorf("expected bitfield message, got %s", d.MessageNames[pm.Id])
	}
	// fmt.Printf("Received %s message with length: %d and payload hex: '%x'\n", d.MessageNames[pm.Id], pm.Length, pm.Payload)

	// Make an interested message and send it
	interestedMessage := d.InterestedMessage()
	// fmt.Printf("Sending interested message to peer: %s\n", conn.RemoteAddr().String())
	// d.LogMessage(interestedMessage.Encode(), false)
	_, err = conn.Write(interestedMessage.Encode())
	if err != nil {
		return nil, fmt.Errorf("error while sending interested message: %v", err)
	}

	// Wait to receive the unchoke message
	size, err = conn.Read(buff)
	if err != nil {
		return nil, fmt.Errorf("error while reading unchoke message: %v", err)
	}
	data = buff[:size]
	// d.LogMessage(data, true)

	// Decode the unchoke message
	decoder = d.NewDecoder(data)
	pm = decoder.DecodePeerMessage()
	if pm.Id != d.UNCHOKE {
		return nil, fmt.Errorf("expected unchoke message, got %s", d.MessageNames[pm.Id])
	}
	// fmt.Printf("Received %s message with length: %d and payload hex: '%x'\n", d.MessageNames[pm.Id], pm.Length, pm.Payload)
	return conn, nil
}

// Create a list of request messages for a piece
func createRequests(pieceLength, numOfBlocks, pieceIndex int) []d.PeerMessage {
	requests := make([]d.PeerMessage, 0)
	for i := 0; i < numOfBlocks; i++ {
		begin := i * BLOCK_LENGTH
		length := BLOCK_LENGTH
		if i == numOfBlocks-1 && pieceLength%BLOCK_LENGTH != 0 {
			// If the piece is the last piece and the piece does not divide evenly into the piece length, adjust the piece length
			length = pieceLength % BLOCK_LENGTH
			// fmt.Printf("Adjusting block length from: %d to: %d\n", BLOCK_LENGTH, pieceLength%BLOCK_LENGTH)
		}
		// fmt.Printf("Adding request message for piece: %d, begin: %d, length: %d\n", pieceIndex, begin, length)
		requests = append(requests, *d.RequestMessage(uint32(pieceIndex), uint32(begin), uint32(length)))
	}
	return requests
}

// Send the request messages to the peer, in chunks of 5 request messages at a time or less if there are fewer than 5 requests left
func sendRequests(requests []d.PeerMessage, conn net.Conn) error {
	count, numOfChainedRequests := 1, 0
	chainedRequests := make([]byte, 0)
	// fmt.Printf("Aiming to send %d request messages\n", len(requests))
	for _, r := range requests {
		chainedRequests = append(chainedRequests, r.Encode()...)
		numOfChainedRequests++
		// Send the request messages in chunks of 5, or less if there are fewer than 5 requests left
		if count%5 == 0 || count+1 > len(requests) {
			// fmt.Printf("Sending %d chained request messages\n", numOfChainedRequests)
			_, err := conn.Write(chainedRequests)
			if err != nil {
				return fmt.Errorf("error while sending request messages: %v", err)
			}
			chainedRequests = make([]byte, 0)
			numOfChainedRequests = 0
		}
		count++
	}
	// fmt.Printf("Sent %d request messages over %d expected\n", count, len(requests))
	return nil
}

// Receive the piece from the peer
// Each piece block has an expected size of 16 kiB + 13 bytes
// Because we're using TCP, we can't guarantee that the piece message will arrive in one piece
func receivePiece(conn net.Conn, pieceLength int) ([]byte, error) {
	stream := make([]byte, 0)
	for len(stream) < pieceLength { // While we haven't received the whole piece
		buff := make([]byte, PIECE_MESSAGE_LENGTH)
		size, err := conn.Read(buff)
		if err != nil {
			return nil, fmt.Errorf("error while reading piece message: %v", err)
		}
		data := buff[:size]
		stream = append(stream, data...)
		// fmt.Printf("Read %d bytes over %d\n", len(stream), pieceLength)
	}
	fmt.Println("Received the whole piece")
	return stream, nil
}

// Chunk the piece message into blocks of 16 kiB + 13 bytes
// If the piece is the last piece and the piece does not divide evenly into the piece length, adjust the piece length
func chunkPieceStream(pieceStream []byte, pieceLength int, isLastPiece, isSameLength bool) [][]byte {
	blocks := make([][]byte, 0)
	// fmt.Printf("Piece stream length: %d\n", len(pieceStream))
	for i := 0; i < len(pieceStream); i += PIECE_MESSAGE_LENGTH {
		var block []byte
		if isLastPiece && !isSameLength && i+PIECE_MESSAGE_LENGTH > pieceLength {
			// fmt.Printf("Adjusting piece message length to %d\n", len(pieceStream)%PIECE_MESSAGE_LENGTH)
			// fmt.Printf("Creating block from %d to %d\n", i, i+len(pieceStream)%PIECE_MESSAGE_LENGTH)
			block = pieceStream[i : i+len(pieceStream)%PIECE_MESSAGE_LENGTH]
		} else {
			block = pieceStream[i : i+PIECE_MESSAGE_LENGTH]
			// fmt.Printf("Creating block from %d to %d\n", i, i+PIECE_MESSAGE_LENGTH)
		}
		blocks = append(blocks, block)
	}
	return blocks
}
