package command

import (
	"fmt"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

// Info prints the information of a torrent file, error if the file cannot be open or decoded
func Info(file string) string {
	torrent, err := OpenTorrentFile(file)
	if err != nil {
		return fmt.Sprintf("error opening torrent file: %v", err)
	}
	return torrent.String()
}

// Open and decode a torrent file
func OpenTorrentFile(file string) (*decoder.TorrentFile, error) {
	fileContent, err := utils.ReadFile(file)
	if err != nil {
		return nil, err
	}
	torrent, _, err := decoder.DecodeTorrentFile(fileContent.String())
	if err != nil {
		return nil, err
	}
	return torrent, nil
}
