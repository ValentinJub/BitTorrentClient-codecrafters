package command

import (
	d "github.com/codecrafters-io/bittorrent-starter-go/decoder"
)

func MagnetParse(magnetLink string) (*d.MagnetLink, error) {
	// Parse the magnet link
	// Example magnet link:
	// magnet:?xt=urn:btih:ad42ce8109f54c99613ce38f9b4d87e70f24a165&dn=magnet1.gif&tr=http%3A%2F%2Fbittorrent-test-tracker.codecrafters.io%2Fannounce
	magnet, err := d.ParseMagnetLink(magnetLink)
	if err != nil {
		return nil, err
	}
	return magnet, nil
}
