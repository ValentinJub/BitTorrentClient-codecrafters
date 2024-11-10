package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/command"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mybittorrent <command> <args>")
		os.Exit(1)
	}
	cHandler := &command.CommandHandlerImpl{}
	cHandler.HandleCommand(os.Args[1], os.Args[2:])
}
