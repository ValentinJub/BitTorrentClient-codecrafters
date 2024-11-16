package command

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/bittorrent-starter-go/encoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

// Handshake performs a handshake with a peer, given the torrent info hash and the peer address
func Handshake(torrentInfoHash, peerAddr string, extended bool) (net.Conn, error) {
	fmt.Println("Handshaking with peer: " + peerAddr)
	// Generate random peer ID
	peerID := utils.GeneratePeerID()
	handshakeMessage := []byte{}
	if !extended {
		handshakeMessage = encoder.MakeHandshakeMessage(torrentInfoHash, peerID, false)
	} else {
		handshakeMessage = encoder.MakeHandshakeMessage(torrentInfoHash, peerID, true)
	}
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to peer: " + err.Error())
	}
	_, err = conn.Write(handshakeMessage)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Error sending handshake to peer: " + err.Error())
	}
	_, err = conn.Read(handshakeMessage)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Error receiving handshake message from peer: " + err.Error())
	}
	fmt.Printf("Peer ID: %x\n", handshakeMessage[48:])
	return conn, nil
}
