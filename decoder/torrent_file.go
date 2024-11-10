package decoder

import (
	"fmt"

	"github.com/codecrafters-io/bittorrent-starter-go/encoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

type TorrentFile struct {
	Announce    string
	Length      int
	InfoHash    string
	PieceLength int
	PieceHashes []string
}

func NewTorrentFile(announce string, length int, infoHash string, pieceLength int, pieceHashes []string) *TorrentFile {
	return &TorrentFile{
		Announce:    announce,
		Length:      length,
		InfoHash:    infoHash,
		PieceLength: pieceLength,
		PieceHashes: pieceHashes,
	}
}

func (t *TorrentFile) String() string {
	hashes := ""
	for _, h := range t.PieceHashes {
		hashes += h + "\n"
	}
	return fmt.Sprintf("Tracker URL: %s\nLength: %d\nInfo Hash: %x\nPiece Length: %d\nPiece Hashes:\n%s", t.Announce, t.Length, t.InfoHash, t.PieceLength, hashes)
}

// A Torrent file is a bencoded dictionary containing information about the torrent
func DecodeTorrentFile(fileContent string) (t *TorrentFile, bytesRead int, err error) {
	decoded, _, err := decodeDictionary(fileContent)
	if err != nil {
		return nil, 0, fmt.Errorf("error decoding torrent file: %v", err)
	}
	URL, ok := decoded["announce"].(string)
	if !ok {
		fmt.Println("No tracker URL found")
		return
	}
	info, ok := decoded["info"].(map[string]interface{})
	if !ok {
		fmt.Println("No info found")
		return
	}
	infoBencoded, err := encoder.EncodeBencode(info)
	if err != nil {
		fmt.Println(err)
		return
	}
	infoHash := utils.SHA1Hash([]byte(infoBencoded))
	length, ok := info["length"].(int)
	if !ok {
		fmt.Println("No length found")
		return
	}
	pieceLength, ok := info["piece length"].(int)
	if !ok {
		fmt.Println("No piece length found")
		return
	}
	pieces, ok := info["pieces"].(string)
	if !ok {
		fmt.Println("No pieces found")
		return
	}
	pieceHashes := make([]string, 0)
	for i := 0; i < len(pieces); i += 20 {
		if i+20 > len(pieces) {
			fmt.Println("Invalid pieces length")
			return
		}
		pieceHashes = append(pieceHashes, fmt.Sprintf("%x", pieces[i:i+20]))
	}
	return NewTorrentFile(URL, length, infoHash, pieceLength, pieceHashes), bytesRead, nil
}
