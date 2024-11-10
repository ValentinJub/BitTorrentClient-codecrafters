package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
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
	case "info":
		fileContent, err := utils.ReadFile(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		decoded, _, err := decoder.DecodeTorrentFile(fileContent.String())
		if err != nil {
			fmt.Println(err)
			return
		}
		URL, ok := decoded["announce"].(string)
		if !ok {
			fmt.Println("No tracker URL found")
			return
		}
		length, ok := decoded["info"].(map[string]interface{})["length"].(int)
		if !ok {
			fmt.Println("No length found")
			return
		}
		fmt.Printf("Tracker URL: %s\nLength: %d\n", URL, length)
	default:
		fmt.Println("Unknown command: " + command)
	}
}
