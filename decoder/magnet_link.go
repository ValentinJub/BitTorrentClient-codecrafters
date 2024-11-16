package decoder

import (
	"fmt"
	"net/url"
	"strings"
)

type MagnetLink struct {
	InfoHash    string
	DisplayName string
	Tracker     string
}

func NewMagnetLink(infoHash string, displayName string, tracker string) *MagnetLink {
	return &MagnetLink{
		InfoHash:    infoHash,
		DisplayName: displayName,
		Tracker:     tracker,
	}
}

func (m *MagnetLink) String() string {
	return fmt.Sprintf("MagnetLink{InfoHash: %s,\n DisplayName: %s,\n Tracker: %s\n}", m.InfoHash, m.DisplayName, m.Tracker)
}

func ParseMagnetLink(magnetLink string) (*MagnetLink, error) {
	// Parse the magnet link
	// Example magnet link:
	// magnet:?xt=urn:btih:ad42ce8109f54c99613ce38f9b4d87e70f24a165&dn=magnet1.gif&tr=http%3A%2F%2Fbittorrent-test-tracker.codecrafters.io%2Fannounce
	// These are the query parameters in a magnet link:

	// xt: urn:btih: followed by the 40-char hex-encoded info hash (example: urn:btih:ad42ce8109f54c99613ce38f9b4d87e70f24a165)
	// dn: The name of the file to be downloaded (example: magnet1.gif)
	// tr: The tracker URL (example: http://bittorrent-test-tracker.codecrafters.io/announce)
	const (
		xt = "xt=urn:btih:"
		dn = "dn="
		tr = "tr="
	)

	n := strings.Index(magnetLink, xt)
	if n == -1 {
		return nil, fmt.Errorf("error while parsing magnet link, could not find info hash")
	}
	infoHash := magnetLink[n+len(xt) : n+len(xt)+40]

	displayName := ""
	n = strings.Index(magnetLink, dn)
	if n == -1 {
		fmt.Printf("Did not find displayName in magnet link\n")
	} else {
		for i := n + len(dn); i < len(magnetLink); i++ {
			if magnetLink[i] == '&' {
				break
			}
			displayName += string(magnetLink[i])
		}
	}

	n = strings.Index(magnetLink, tr)
	tracker := ""
	if n == -1 {
		fmt.Printf("Did not find tracker in magnet link\n")
	} else {
		for i := n + len(tr); i < len(magnetLink); i++ {
			if magnetLink[i] == '&' {
				break
			}
			tracker += string(magnetLink[i])
		}
		// Decode the tracker URL
		tracker, _ = decodeURL(tracker)
	}

	return NewMagnetLink(infoHash, displayName, tracker), nil
}

func decodeURL(addr string) (string, error) {
	// Decode the URL
	// Example: http%3A%2F%2Fbittorrent-test-tracker.codecrafters.io%2Fannounce
	// should be decoded to: http://bittorrent-test-tracker.codecrafters.io/announce
	decodedURL, err := url.QueryUnescape(addr)
	if err != nil {
		return "", fmt.Errorf("error while decoding URL: %v", err)
	}
	return decodedURL, nil
}
