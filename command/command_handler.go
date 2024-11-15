package command

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

type CommandHanlder interface {
	HandleCommand(command string, args []string) (string, error)
}

type CommandHandlerImpl struct {
}

func (c *CommandHandlerImpl) HandleCommand(command string, args []string) {
	if len(args) < 1 || command == "" {
		fmt.Println("Usage: mybittorrent <command> <args>")
		return
	}
	switch command {
	// $ ./your_bittorrent.sh decode <bencoded_string>
	// example:
	// $ ./your_bittorrent.sh decode d3:foo3:bar5:helloi52ee
	case "decode":
		Decode(args[0])
	// $ ./your_bittorrent.sh download_piece -o <output_dir> <torrent.file> <piece_index>
	// example:
	// $ ./your_bittorrent.sh download_piece -o output sample.torrent 0
	case "download_piece":
		if len(args) < 4 {
			fmt.Println("Usage: mybittorrent download_piece -o <output_dir> <torrent.file> <piece_index>")
			return
		}
		// for debugging
		// fmt.Printf("args %v\n", args)
		pieceIndex, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Invalid piece index: ", err)
			return
		}
		// Get the torrent file information
		torrentFile := args[2]
		torrent, err := OpenTorrentFile(torrentFile)
		if err != nil {
			fmt.Println("Error while opening torrent file: ", err)
			return
		}
		// Get the list of peers from the tracker
		peers, err := Peers(torrent.Announce, torrent.InfoHash, torrent.Length)
		if err != nil {
			return
		}
		outputDir := args[1]
		piece, err := DownloadPiece(peers[0], torrent, pieceIndex)
		if err != nil {
			fmt.Println("Error while downloading piece: ", err)
			return
		}
		// Write the piece to the output file
		err = utils.WriteFile(outputDir, piece)
		if err != nil {
			fmt.Println("Error while writing piece to file: ", err)
		}
	// $ ./your_bittorrent.sh handshake sample.torrent <peer_ip>:<peer_port>
	case "handshake":
		if len(args) < 2 {
			fmt.Println("Usage: mybittorrent handshake <torrent.file> <peer_ip:port>")
			return
		}
		torrentFileName, peerAddr := args[0], args[1]
		torrentFile, err := OpenTorrentFile(torrentFileName)
		if err != nil {
			fmt.Println("Error while opening torrent file: ", err)
			return
		}
		Handshake(torrentFile.InfoHash, peerAddr)
	// $ ./your_bittorrent.sh info sample.torrent
	case "info":
		if len(args) < 1 {
			fmt.Println("Usage: mybittorrent info <torrent.file>")
			return
		}
		torrentFile := args[0]
		fmt.Print(Info(torrentFile))
	// $ ./your_bittorrent.sh peers sample.torrent
	case "peers":
		if len(args) < 1 {
			fmt.Println("Usage: mybittorrent peers <torrent.file>")
			return
		}
		torrentFileName := args[0]
		torrent, err := OpenTorrentFile(torrentFileName)
		if err != nil {
			fmt.Println("Error while opening torrent file: ", err)
			return
		}
		peers, err := Peers(torrent.Announce, torrent.InfoHash, torrent.Length)
		if err != nil {
			fmt.Println("Error while getting peers: ", err)
			return
		}
		// Print the peers
		for _, peer := range peers {
			fmt.Println(peer)
		}
	default:
		fmt.Println("Unknown command: " + command)
	}
}
