package command

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/encoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

func Handshake(torrent, peerAddr string) {
	// Get the sha1 info hash from the torrent file
	fileContent, err := utils.ReadFile(torrent)
	if err != nil {
		fmt.Println("Error reading torrent file: " + err.Error())
		return
	}
	t, _, err := decoder.DecodeTorrentFile(fileContent.String())
	if err != nil {
		fmt.Println("Error decoding torrent file: " + err.Error())
		return
	}
	infoHash := t.InfoHash
	// Generate random peer ID
	peerID := utils.GeneratePeerID()
	fmt.Println("Handshaking with peer: " + peerAddr)
	// Handshake with the peer
	// Send the handshake message to the peer
	handshakeMessage := encoder.MakeHandshakeMessage(infoHash, peerID)
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		fmt.Println("Error connecting to peer: " + err.Error())
		return
	}
	defer conn.Close()
	_, err = conn.Write(handshakeMessage)
	if err != nil {
		fmt.Println("Error sending handshake message to peer: " + err.Error())
		return
	}
	_, err = conn.Read(handshakeMessage)
	if err != nil {
		fmt.Println("Error receiving handshake message from peer: " + err.Error())
		return
	}
	fmt.Printf("Peer ID: %x\n", handshakeMessage[48:])
}
