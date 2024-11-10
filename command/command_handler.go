package command

import (
	"fmt"
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
	case "decode":
		Decode(args[0])
	case "handshake":
		if len(args) < 2 {
			fmt.Println("Usage: mybittorrent handshake <torrent.file> <peer_ip:port>")
			return
		}
		Handshake(args[0], args[1])
	case "info":
		Info(args[0])
	case "peers":
		Peers(args[0])
	default:
		fmt.Println("Unknown command: " + command)
	}
}
