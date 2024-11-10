package command

import (
	"fmt"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

func Info(file string) {
	fileContent, err := utils.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	torrent, _, err := decoder.DecodeTorrentFile(fileContent.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(torrent.String())
}
