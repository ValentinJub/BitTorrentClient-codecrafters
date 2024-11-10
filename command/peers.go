package command

import (
	"fmt"
	"net/url"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/netclient"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

func Peers(file string) {
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
	// Make a GET request to the tracker URL
	client := &netclient.Client{
		RemoteURL: torrent.Announce,
	}
	queryParameters := fmt.Sprintf("?info_hash=%s&peer_id=%s&port=%d&uploaded=%d&downloaded=%d&left=%d&compact=1", url.QueryEscape(torrent.InfoHash), "MyCustomIDValentin?!", 6881, 0, 0, torrent.Length)
	req, err := client.CreateRequest("GET", queryParameters, nil)
	if err != nil {
		fmt.Println("Error while creating request: ", err)
		return
	}
	resp, err := client.MakeRequest(req)
	if err != nil {
		fmt.Println("Error while making request: ", err)
		return
	}
	// Decode the response
	decoded, _, err := decoder.DecodeBencode(string(resp))
	if err != nil {
		fmt.Println(err)
		return
	}
	switch decoded.(type) {
	case map[string]interface{}:
		peers := decoded.(map[string]interface{})["peers"].(string)
		// Parse the peers string
		peersList, err := ParsePeers(peers)
		if err != nil {
			fmt.Println("Error while parsing peers: ", err)
			return
		}
		// Print the peers
		for _, peer := range peersList {
			fmt.Println(peer)
		}
	default:
		fmt.Println("Unsupported type")
	}

}

func ParsePeers(peers string) ([]string, error) {
	// The peers string is a string of 6 bytes for each peer
	// The first 4 bytes are the IP address and the last 2 bytes are the port number
	peersList := []string{}
	for i := 0; i < len(peers); i += 6 {
		if i+6 > len(peers) {
			break
		}
		peer := fmt.Sprintf("%d.%d.%d.%d:%d", peers[i], peers[i+1], peers[i+2], peers[i+3], int(peers[i+4])<<8|int(peers[i+5]))
		peersList = append(peersList, peer)
	}
	return peersList, nil
}
