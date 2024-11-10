package netclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	RemoteURL string
}

// Create a new request with the given method, url and body
func (c *Client) CreateRequest(method, urlpart string, bodyData *bytes.Buffer) (*http.Request, error) {
	var req *http.Request
	var err error

	if method == "GET" || bodyData == nil {
		req, err = http.NewRequest(method, c.RemoteURL+urlpart, nil)
		// Log the request URL for debugging
		// fmt.Printf("Request URL: %s\n", c.RemoteURL+urlpart)
	} else {
		req, err = http.NewRequest(method, c.RemoteURL+urlpart, bodyData)
	}

	if err != nil {
		return nil, err
	}
	return req, nil
}

// Make a request and decode the response into the target interface
func (c *Client) MakeRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while making request: ", err)
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading the response: ", err)
		return []byte{}, err
	}
	// Log the response body for debugging
	// fmt.Printf("Length: %d - Response body: %s\n", len(body), string(body))

	return body, nil
}
