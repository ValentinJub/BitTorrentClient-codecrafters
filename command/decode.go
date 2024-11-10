package command

import (
	"encoding/json"
	"fmt"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
)

func Decode(value string) {
	decoded, _, err := decoder.DecodeBencode(value)
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
}
