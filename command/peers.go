package command

import (
	"fmt"
	"net/url"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/netclient"
)

// Peers gets the list of peers from a torrent file
// It sends a request to the tracker and gets the list of peers
// The list of peers is a string of 6 bytes for each peer
// The first 4 bytes are the IP address and the last 2 bytes are the port number
// The function returns a list of strings with the IP address and port number of each peer
func Peers(announceURL, torrentInfoHash string, torrentLength int) ([]string, error) {
	client := &netclient.Client{
		RemoteURL: announceURL,
	}
	queryParameters := fmt.Sprintf("?info_hash=%s&peer_id=%s&port=%d&uploaded=%d&downloaded=%d&left=%d&compact=1", url.QueryEscape(torrentInfoHash), "MyCustomIDValentin?!", 6881, 0, 0, torrentLength)
	req, err := client.CreateRequest("GET", queryParameters, nil)
	if err != nil {
		return []string{}, fmt.Errorf("error while creating request: %s", err.Error())
	}
	resp, err := client.MakeRequest(req)
	if err != nil {
		return []string{}, fmt.Errorf("error while making request: %s", err.Error())
	}
	// Decode the response
	decoded, _, err := decoder.DecodeBencode(string(resp))
	if err != nil {
		return []string{}, fmt.Errorf("error while decoding the response: %s", err.Error())
	}
	switch decoded.(type) {
	case map[string]interface{}:
		peers := decoded.(map[string]interface{})["peers"].(string)
		// Parse the peers string
		peersList, err := ParsePeers(peers)
		if err != nil {
			return []string{}, fmt.Errorf("error while parsing the peers: %s", err.Error())
		}
		return peersList, nil
	default:
		return []string{}, fmt.Errorf("error unexpected response: %v", decoded)
	}
}

// The peers string is a string of 6 bytes for each peer
// The first 4 bytes are the IP address and the last 2 bytes are the port number
func ParsePeers(peers string) ([]string, error) {
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
