package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
)

type CommandHanlder interface {
	HandleCommand(command string, args []string) (string, error)
}

type CommandHandlerImpl struct {
}

func (c *CommandHandlerImpl) HandleCommand(command string, args []string) {
	switch command {
	case "decode":
		bencodedValue := os.Args[2]

		decoded, _, err := decoder.DecodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		if string(jsonOutput) == "null" {
			fmt.Println("[]")
		} else {
			fmt.Println(string(jsonOutput))
		}
	default:
		fmt.Println("Unknown command: " + command)
	}
}
